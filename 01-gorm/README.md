# 官方

官方文档：https://gorm.io/zh_CN/docs/index.html

扩展：https://github.com/jinzhu/gorm

# 增删改查

https://gorm.io/zh_CN/docs/advanced_query.html

# 链表

https://gorm.io/zh_CN/docs/query.html#Joins

# 事务

https://gorm.io/zh_CN/docs/transactions.html

# Scope

https://gorm.io/zh_CN/docs/scopes.html





# Page

```go
package handler

import "gorm.io/gorm"

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

```

调用page

```go
func (s *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error){
	brandListResponse := proto.BrandListResponse{}

	var brands []model.Brands
    // 调用page
	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}

	var total int64
	global.DB.Model(&model.Brands{}).Count(&total)
	brandListResponse.Total = int32(total)

	var brandResponses []*proto.BrandInfoResponse
	for _, brand := range brands {
		brandResponses = append(brandResponses, &proto.BrandInfoResponse{
			Id:  brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	brandListResponse.Data = brandResponses
	return &brandListResponse, nil
}
```



