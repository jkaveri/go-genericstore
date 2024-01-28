package converter_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goflexstore/converter"
)

type UnMatchUser struct {
	ID   int
	Name int
}

func (e UnMatchUser) GetID() int {
	return e.ID
}

func Test_Converter_ToEntity(t *testing.T) {
	t.Run("should-convert-DTO-to-Entity", func(t *testing.T) {
		converter := converter.NewReflect[User, UserDTO, int](nil)
		now := time.Now()

		dto := UserDTO{
			ID:        1,
			Name:      "name",
			Age:       10,
			Disabled:  sql.NullBool{Bool: false, Valid: true},
			IsAdmin:   &sql.NullBool{Bool: true, Valid: true},
			CreatedAt: sql.NullTime{Valid: true, Time: now},
			Referer: &UserDTO{
				ID:        2,
				Name:      "referer",
				Age:       43,
				Disabled:  sql.NullBool{Bool: false, Valid: true},
				IsAdmin:   &sql.NullBool{Bool: true, Valid: true},
				CreatedAt: sql.NullTime{Valid: true, Time: now},
			},
			Friends: []*UserDTO{
				{
					ID:        3,
					Name:      "friend1",
					Age:       32,
					Disabled:  sql.NullBool{Bool: false, Valid: true},
					IsAdmin:   &sql.NullBool{Bool: true, Valid: true},
					CreatedAt: sql.NullTime{Valid: true, Time: now},
				},
			},
		}

		entity := converter.ToEntity(dto)

		assert.Equal(t,
			User{
				ID:        1,
				Name:      "name",
				Age:       10,
				IsAdmin:   true,
				CreatedAt: now,

				Referer: &User{
					ID:        2,
					Name:      "referer",
					Age:       43,
					IsAdmin:   true,
					CreatedAt: now,
				},
				Friends: []*User{
					{
						ID:        3,
						Name:      "friend1",
						Age:       32,
						Disabled:  false,
						IsAdmin:   true,
						CreatedAt: now,
					},
				},
			},
			entity,
		)
	})

	t.Run("map-from-pointer-type", func(t *testing.T) {
		converter := converter.NewReflect[*User, *UserDTO, int](nil)

		dto := UserDTO{ID: 1, Name: "name", Disabled: sql.NullBool{Bool: false, Valid: true}}

		entity := converter.ToEntity(&dto)

		assert.Equal(t, &User{ID: 1, Name: "name"}, entity)
	})

	t.Run("map-from-pointer-type-nil-val", func(t *testing.T) {
		converter := converter.NewReflect[*User, *UserDTO, int](nil)

		var dto *UserDTO

		entity := converter.ToEntity(dto)

		assert.Equal(t, (*User)(nil), entity)
	})

	t.Run("should-convert-empty-DTO-to-empty-Entity", func(t *testing.T) {
		converter := converter.NewReflect[User, UserDTO, int](nil)
		dto := UserDTO{}

		entity := converter.ToEntity(dto)

		assert.Equal(t, User{}, entity)
	})

	t.Run("should-panic", func(t *testing.T) {
		converter := converter.NewReflect[UnMatchUser, UserDTO, int](nil)
		dto := UserDTO{
			ID:       1,
			Name:     "John",
			Age:      3,
			IsAdmin:  &sql.NullBool{},
			Disabled: sql.NullBool{},
		}

		assert.PanicsWithError(t, "cannot assign src.Name(string) to dst.Name(int)", func() {
			_ = converter.ToEntity(dto)
		})
	})
}

func Test_ToMany(t *testing.T) {
	t.Run("should-convert-DTOs-to-Entities", func(t *testing.T) {
		conv := converter.NewReflect[User, UserDTO, int](nil)
		dtos := []UserDTO{
			{ID: 1, Name: "name1"},
			{ID: 2, Name: "name2"},
		}

		entities := converter.ToMany(dtos, conv.ToEntity)

		assert.Equal(t,
			[]User{
				{ID: 1, Name: "name1"},
				{ID: 2, Name: "name2"},
			},
			entities,
		)
	})

	t.Run("should-convert-empty-DTOs-to-empty-Entities", func(t *testing.T) {
		conv := converter.NewReflect[User, UserDTO, int](nil)
		dtos := []UserDTO{}

		entities := converter.ToMany(dtos, conv.ToEntity)

		assert.Nil(t, entities)
	})
}

func Test_Converter_ToDTO(t *testing.T) {
	t.Run("should-convert-Entity-to-DTO", func(t *testing.T) {
		now := time.Now()
		conv := converter.NewReflect[User, UserDTO, int](nil)

		entity := User{
			ID: 1, Name: "name",
			Age:       10,
			Disabled:  true,
			IsAdmin:   true,
			CreatedAt: now,
		}

		dto := conv.ToDTO(entity)

		assert.Equal(t, UserDTO{
			ID:        1,
			Name:      "name",
			Age:       10,
			IsAdmin:   &sql.NullBool{Valid: true, Bool: true},
			Disabled:  sql.NullBool{Valid: true, Bool: true},
			CreatedAt: sql.NullTime{Valid: true, Time: now},
		}, dto)
	})
}
