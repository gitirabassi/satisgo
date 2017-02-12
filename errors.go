package satisgo

import "fmt"

func handleHeader(header int) error {
	switch header {
	case 200:
		return nil
	case 204:
		return nil
	case 400:
		return fmt.Errorf("%d -- Bad Request – Your request is wrong", header)
	case 401:
		return fmt.Errorf("%d -- Unauthorized – Your API key is wrong", header)
	case 403:
		return fmt.Errorf("%d -- Forbidden – The resource requested is hidden for administrators only", header)
	case 404:
		return fmt.Errorf("%d -- Not Found – The specified resource could not be found", header)
	case 500:
		return fmt.Errorf("%d -- Internal Server Error – We had a problem with our server. Try again later", header)
	case 503:
		return fmt.Errorf("%d -- Service Unavailable – We’re temporarially offline for maintanance. Please try again later", header)
	default:
		return fmt.Errorf("%d -- Header status code not found in the standard error cases given by Satispay", header)
	}
}
