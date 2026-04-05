package service

import (
	"sync"

	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type DashboardService struct {
	incomeRepo  *repository.IncomeRepository
	expenseRepo *repository.ExpenseRepository
	debtRepo    *repository.DebtRepository
}

func NewDashboardService(incomeRepo *repository.IncomeRepository, expenseRepo *repository.ExpenseRepository, debtRepo *repository.DebtRepository) *DashboardService {
	return &DashboardService{incomeRepo: incomeRepo, expenseRepo: expenseRepo, debtRepo: debtRepo}
}

type MonthlyNet struct {
	Month    int    `json:"month"`
	Income   string `json:"income"`
	Expenses string `json:"expenses"`
	Net      string `json:"net"`
}

type CategorySum struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	Total        string `json:"total"`
}

type TypeSum struct {
	Type  string `json:"type"`
	Total string `json:"total"`
}

type IncomeVsExpenses struct {
	TotalIncome   string `json:"total_income"`
	TotalExpenses string `json:"total_expenses"`
}

type DailySum struct {
	Date  string `json:"date"`
	Total string `json:"total"`
}

type DashboardData struct {
	NetSavings       []MonthlyNet     `json:"net_savings"`
	IncomeBreakdown  []CategorySum    `json:"income_breakdown"`
	ExpenseBreakdown []TypeSum        `json:"expense_breakdown"`
	IncomeVsExpenses IncomeVsExpenses `json:"income_vs_expenses"`
	DailyExpenses    []DailySum       `json:"daily_expenses"`
	DailyDebts       []DailySum       `json:"daily_debts"`
}

func (s *DashboardService) GetDashboard(userID uuid.UUID, year int, month *int) (*DashboardData, error) {
	var (
		incomeSums    []repository.MonthSum
		expenseSums   []repository.MonthSum
		debtSums      []repository.MonthSum
		categorySums  []repository.CategorySumRow
		typeSums      []repository.TypeSumRow
		totalIncome   decimal.Decimal
		totalExpenses decimal.Decimal
		dailyExpenses []repository.DaySumRow
		dailyDebts    []repository.DaySumRow
		mu            sync.Mutex
		firstErr      error
		wg            sync.WaitGroup
	)

	setErr := func(err error) {
		mu.Lock()
		if firstErr == nil {
			firstErr = err
		}
		mu.Unlock()
	}

	wg.Add(9)
	go func() { defer wg.Done(); r, e := s.incomeRepo.SumByMonth(userID, year); if e != nil { setErr(e) } else { incomeSums = r } }()
	go func() { defer wg.Done(); r, e := s.expenseRepo.SumByMonth(userID, year); if e != nil { setErr(e) } else { expenseSums = r } }()
	go func() { defer wg.Done(); r, e := s.debtRepo.SumByMonth(userID, year); if e != nil { setErr(e) } else { debtSums = r } }()
	go func() { defer wg.Done(); r, e := s.incomeRepo.SumByCategory(userID, year); if e != nil { setErr(e) } else { categorySums = r } }()
	go func() { defer wg.Done(); r, e := s.expenseRepo.SumByType(userID, year); if e != nil { setErr(e) } else { typeSums = r } }()
	go func() { defer wg.Done(); r, e := s.incomeRepo.TotalByYear(userID, year); if e != nil { setErr(e) } else { totalIncome = r } }()
	go func() { defer wg.Done(); r, e := s.expenseRepo.TotalByYear(userID, year); if e != nil { setErr(e) } else { totalExpenses = r } }()
	go func() { defer wg.Done(); r, e := s.expenseRepo.SumByDay(userID, year, month); if e != nil { setErr(e) } else { dailyExpenses = r } }()
	go func() { defer wg.Done(); r, e := s.debtRepo.SumByDay(userID, year, month); if e != nil { setErr(e) } else { dailyDebts = r } }()
	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	// Build monthly net savings
	incomeMap := make(map[int]decimal.Decimal)
	for _, s := range incomeSums {
		incomeMap[s.Month] = s.Total
	}
	expenseMap := make(map[int]decimal.Decimal)
	for _, s := range expenseSums {
		expenseMap[s.Month] = s.Total
	}
	debtMap := make(map[int]decimal.Decimal)
	for _, s := range debtSums {
		debtMap[s.Month] = s.Total
	}

	netSavings := make([]MonthlyNet, 12)
	for m := 1; m <= 12; m++ {
		inc := incomeMap[m]
		exp := expenseMap[m].Add(debtMap[m])
		netSavings[m-1] = MonthlyNet{
			Month:    m,
			Income:   inc.String(),
			Expenses: exp.String(),
			Net:      inc.Sub(exp).String(),
		}
	}

	incomeBreakdown := make([]CategorySum, len(categorySums))
	for i, cs := range categorySums {
		incomeBreakdown[i] = CategorySum{CategoryID: cs.CategoryID, CategoryName: cs.CategoryName, Total: cs.Total.String()}
	}

	expenseBreakdown := make([]TypeSum, len(typeSums))
	for i, ts := range typeSums {
		expenseBreakdown[i] = TypeSum{Type: ts.Type, Total: ts.Total.String()}
	}

	dailyExpOut := make([]DailySum, len(dailyExpenses))
	for i, d := range dailyExpenses {
		dailyExpOut[i] = DailySum{Date: d.Date, Total: d.Total.String()}
	}
	dailyDebtOut := make([]DailySum, len(dailyDebts))
	for i, d := range dailyDebts {
		dailyDebtOut[i] = DailySum{Date: d.Date, Total: d.Total.String()}
	}

	return &DashboardData{
		NetSavings:       netSavings,
		IncomeBreakdown:  incomeBreakdown,
		ExpenseBreakdown: expenseBreakdown,
		IncomeVsExpenses: IncomeVsExpenses{TotalIncome: totalIncome.String(), TotalExpenses: totalExpenses.String()},
		DailyExpenses:    dailyExpOut,
		DailyDebts:       dailyDebtOut,
	}, nil
}
