package handlers

import (
	"jobportalapi/internal/auth"
	"jobportalapi/internal/services"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAPI(t *testing.T) {
	type args struct {
		a *auth.Auth
		c *services.Conn
	}
	tests := []struct {
		name string
		args args
		want *gin.Engine
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := API(tt.args.a, tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("API() = %v, want %v", got, tt.want)
			}
		})
	}
}
