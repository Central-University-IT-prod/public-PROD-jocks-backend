package service

import "github.com/jackc/pgx/v5/pgxpool"

type Services struct {
	Bs               *BusinessService
	Ps               *ProductService
	Cs               *CheckService
	Fs               *FormService
	ClientService    *ClientService
	Rs               *ReviewService
	PromocodeService *PromocodeService
}

func NewService(pool *pgxpool.Pool) (*Services, error) {
	bs, err := NewBusinessService(pool)
	if err != nil {
		return nil, err
	}

	ps, err := NewProductService(pool)
	if err != nil {
		return nil, err
	}
	cs, err := NewCheckService(pool)
	if err != nil {
		return nil, err
	}
	fs, err := NewFormService(pool)
	if err != nil {
		return nil, err
	}

	cliS, err := NewClientService(pool)
	if err != nil {
		return nil, err
	}
	rs, err := NewReviewService(pool)
	if err != nil {
		return nil, err
	}

	promoService, err := NewPromocodeService(pool)
	if err != nil {
		return nil, err
	}

	return &Services{Bs: bs, Ps: ps, Cs: cs, Fs: fs, ClientService: cliS, Rs: rs, PromocodeService: promoService}, nil
}
