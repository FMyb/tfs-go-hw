package orders

import (
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/domain"
	"testing"
)

func TestSendOrder(t *testing.T) {
	type args struct {
		symbol        string
		size          uint64
		side          string
		publicApiKey  string
		privateApiKey string
	}
	tests := []struct {
		name    string
		args    args
		want    domain.ResponseOrder
		wantErr bool
	}{
		{
			name: "send order",
			args: args{
				symbol:        "PI_XBTUSD",
				size:          1,
				side:          "buy",
				publicApiKey:  "51nxJvodlps7HPcJ5Y1Z0PFJ/+ypHkbaoaULia82rqPiJX2vNnU+RmGn",
				privateApiKey: "RRzV+Hdbg5zCFKLhq1B3YYAVtZozpIcXe8LmEyIaM4EgpM6bPm7NJxNMhrbZKqoC5J9coq0YEm1kMkAS4ZRm3BVh",
			},
			want: domain.KrakenResponseOrder{ // TODO Mock
				Kresult: domain.SUCCESSES,
				Kerror:  "",
			},
			wantErr: false,
		},
		{
			name: "fail order",
			args: args{
				symbol:        "",
				size:          0,
				side:          "",
				publicApiKey:  "",
				privateApiKey: "",
			},
			want: domain.KrakenResponseOrder{
				Kresult: domain.ERROR,
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDefaultKrakenOrders().SendOrder(
				tt.args.symbol,
				tt.args.size,
				tt.args.side,
				tt.args.publicApiKey,
				tt.args.privateApiKey,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Result() != tt.want.Result() {
				t.Errorf("SendOrder() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateAuthent(t *testing.T) {
	type args struct {
		PostData     string
		endPointPath string
		apiKey       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "generateAuthent",
			args: args{
				PostData:     "orderType=mkt&side=buy&size=1&symbol=PI_XBTUSD",
				endPointPath: "api/v3/sendorder",
				apiKey:       "RRzV+Hdbg5zCFKLhq1B3YYAVtZozpIcXe8LmEyIaM4EgpM6bPm7NJxNMhrbZKqoC5J9coq0YEm1kMkAS4ZRm3BVh",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := generateAuthent(tt.args.PostData, tt.args.endPointPath, tt.args.apiKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateAuthent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
