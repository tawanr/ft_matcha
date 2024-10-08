package models

import (
	"testing"

	"github.com/tawanr/ft_matcha/internal/assert"
)

func TestProfile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	db := newTestDB(t)
	p := &ProfileModel{DB: db}
	t.Run("Insert", func(t *testing.T) {
		profile, err := p.Insert(1, 0, []int{0}, "test bio", 18)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, profile.UserID, int64(1))
	})

	t.Run("Get", func(t *testing.T) {
		profile, err := p.Get(1)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, profile.UserID, int64(1))
		assert.Equal(t, profile.Bio, "test bio")
		assert.Equal(t, profile.Gender, 0)
		assert.Equal(t, profile.Age, 18)
	})

	t.Run("Update", func(t *testing.T) {
		profile, err := p.Update(1, 1, []int{1}, "updated bio", 20)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, profile.UserID, int64(1))
		assert.Equal(t, profile.Bio, "updated bio")
		assert.Equal(t, profile.Gender, 1)
		assert.Equal(t, profile.Age, 20)
	})
}
