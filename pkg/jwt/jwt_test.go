package jwt_test

import (
	"go/avito-test/pkg/jwt"
	"reflect"
	"testing"
)

func TestJwt_Parse(t *testing.T) {
	type fields struct {
		Secret string
	}
	type args struct {
		tokenString string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *jwt.JwtData
		want1  bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jwt.Jwt{
				Secret: tt.fields.Secret,
			}
			got, got1 := j.Parse(tt.args.tokenString)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Parse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestJwt_Sign(t *testing.T) {
	const testUserName = "bob"
	jwtService := jwt.NewJwt("5cv2zNCN5zLZJTp4uytfekjLSDkw3ciAI+V7CGT/f9A=")
	token, err := jwtService.Sign(jwt.JwtData{
		Username: testUserName,
	})
	if err != nil {
		t.Fatal(err)
	}
	data, isValid := jwtService.Parse(token)
	if !isValid {
		t.Fatal("Token is not valid")
	}
	if data.Username != testUserName {
		t.Fatalf("Username %s is not equal to %s ", data.Username, testUserName)
	}

}

func TestNewJwt(t *testing.T) {
	type args struct {
		secret string
	}
	tests := []struct {
		name string
		args args
		want *jwt.Jwt
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jwt.NewJwt(tt.args.secret); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJwt() = %v, want %v", got, tt.want)
			}
		})
	}
}
