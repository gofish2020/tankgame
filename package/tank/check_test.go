package tank

import "testing"

func Test_checkCollision1(t *testing.T) {
	type args struct {
		point    Point
		vertices []Point
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				point: Point{
					X: 1,
					Y: 2,
				},
				vertices: []Point{{2, 1}, {2, 2}, {3, 2}, {3, 1}},
			},
			want: false,
		},
		{
			name: "test",
			args: args{
				point: Point{
					X: 2.5,
					Y: 1.5,
				},
				vertices: []Point{{2, 1}, {2, 2}, {3, 2}, {3, 1}},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkCollision1(tt.args.point, tt.args.vertices); got != tt.want {
				t.Errorf("checkCollision1() = %v, want %v", got, tt.want)
			}
		})
	}
}
