package prices

import "testing"

func TestGetPrice(t *testing.T) {
	price := GetPrice("2018-06-01", "BCH","/home/sean/.config/cgtcalc/pricedb")
	t.Logf("2018-07-01-BCH = %.2f", price)
	if price != 1325.145432 {
		t.Errorf("Prices was incorrect, got %.2f, want %.2f",price, 10.0)

	}
}
