package api_customer

import (
	"net/http"
	"server/service"
	"server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QueryBusinessResp struct {
	Items []QueryBusinessRespItem `json:"items"`
}

type QueryBusinessRespItem struct {
	Id               int     `json:"id"`
	Name             string  `json:"name"`
	DescriptionShort string  `json:"description_short"`
	Rating           float64 `json:"rating"`
	RatingCount      int     `json:"rating_count"`
	Image            string  `json:"image"`
	Address          Address `json:"address"`
}

type GetBusinessResp struct {
	Id               int               `json:"id"`
	Name             string            `json:"name"`
	DescriptionShort string            `json:"description_short"`
	Rating           []FormParam       `json:"rating"`
	RatingAverage    float64           `json:"rating_average"`
	RatingCount      int               `json:"rating_count"`
	Image            string            `json:"image"`
	Addresses        []Address         `json:"addresses"`
	Items            []GetBusinessItem `json:"items"`
}

type GetBusinessItem struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Rating      float64 `json:"rating"`
	RatingCount int     `json:"rating_count"`
}

type FormParam struct {
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	Rating       float64 `json:"rating"`
	ReviewsCount int     `json:"reviews_count"`
}

type Address struct {
	Full string  `json:"full"`
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

func getQueryOpts(ctx *gin.Context) service.QueryOpts {
	opts := service.QueryOpts{}
	opts.Query = ctx.Query("query")

	limitStr := ctx.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 15
	}

	offsetStr := ctx.Query("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	opts.Limit = limit
	opts.Offset = offset

	latStr := ctx.Query("lat")
	if latStr == "" {
		return opts
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return opts
	}

	lngStr := ctx.Query("long")
	if latStr == "" {
		return opts
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		return opts
	}

	opts.SearchOpts = &service.SearchOpts{Lat: lat, Lng: lng, Radius: 5}
	return opts
}

func QueryBusiness(s *service.Services) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		opts := getQueryOpts(ctx)

		res, err := s.Bs.QueryBusinesses(opts)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		resp := QueryBusinessResp{}
		items := []QueryBusinessRespItem{}

		for _, business := range res {
			rating, raitingCount, err := s.Rs.GetRatingForBusiness(business.Business.Id)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
				return
			}

			items = append(items, QueryBusinessRespItem{
				Id:               business.Business.Id,
				Name:             business.Business.Name,
				DescriptionShort: business.Business.Description,
				Rating:           rating,
				RatingCount:      raitingCount,
				Image:            "https://img.freepik.com/premium-vector/default-avatar-profile-icon-social-media-user-image-gray-avatar-icon-blank-profile-silhouette-vector-illustration_561158-3383.jpg",
				Address: Address{
					Full: business.Address.Address,
					Lat:  business.Coords.Lat,
					Long: business.Coords.Lng,
				},
			})
		}
		resp.Items = items

		ctx.JSON(http.StatusOK, resp)
	}
}

func GetBuisnessDetailsById(s *service.Services) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
			return
		}

		business, err := s.Bs.GetById(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
			return
		}

		addresses, err := s.Bs.GetAddressesFor(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
			return
		}

		rating, err := s.Rs.GetBusinessFormsRating(business.Id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
			return
		}

		ratingAverage, ratingCount, err := s.Rs.GetRatingForBusiness(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
			return
		}

		products, err := s.Rs.GetRatingForProductsOf(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
			return
		}

		resp := GetBusinessResp{
			Id:               id,
			Name:             business.Name,
			DescriptionShort: business.Description,
			Rating: utils.Map(rating, func(item service.ParamWithRating) FormParam {
				return FormParam{
					Name:         item.Name,
					Type:         item.Type,
					Rating:       item.Rating,
					ReviewsCount: item.RatingCount,
				}
			}),
			RatingAverage: ratingAverage,
			RatingCount:   ratingCount,
			Image:         "https://img.freepik.com/premium-vector/default-avatar-profile-icon-social-media-user-image-gray-avatar-icon-blank-profile-silhouette-vector-illustration_561158-3383.jpg",
			Addresses: utils.Map(addresses, func(item service.Address) Address {
				return Address{
					Full: item.Address,
					Lat:  item.Coords.Lat,
					Long: item.Coords.Lng,
				}
			}),
			Items: utils.Map(products, func(item service.ProductsWithRating) GetBusinessItem {
				return GetBusinessItem{
					Id:          item.Id,
					Name:        item.Name,
					Rating:      item.Rating,
					RatingCount: item.RatingCount,
				}
			}),
		}

		ctx.JSON(http.StatusOK, resp)
	}
}
