package main

import "fmt"

func originalDiscountAll(items []float64) []float64 {
	for i := range items {
		items[i] *= 0.9
	}
	return items
}

func discountAll(items []float64) []float64 {
	result := make([]float64, len(items))
	for i, price := range items {
		result[i] = price * 0.9
	}
	return result
}

func main() {
	originalPrices := []float64{100, 200, 300, 400}
	originalWeekend := originalPrices[0:2]
	originalDiscounted := originalDiscountAll(originalWeekend)
	fmt.Println("original aliasing result:")
	fmt.Println(originalPrices)     // [90 180 300 400]
	fmt.Println(originalWeekend)    // [90 180]
	fmt.Println(originalDiscounted) // [90 180]

	prices := []float64{100, 200, 300, 400}
	weekend := prices[0:2]
	discounted := discountAll(weekend)

	fmt.Println("fixed result:")
	fmt.Println(prices)     // [100 200 300 400]
	fmt.Println(weekend)    // [100 200]
	fmt.Println(discounted) // [90 180]

	a := prices[0:2]
	a = append(a, 999)
	fmt.Println(prices) // [100 200 999 400]

	b := prices[0:2:2]
	b = append(b, 888)
	fmt.Println(prices) // [100 200 999 400] (ไม่ถูก append ทับ)
	fmt.Println(b)      // [100 200 888]
}
