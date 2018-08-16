/*
 * Dummy Service Provider generated using 'swaggercodegen' that has two resources 'cdns' and 'lbs' which are terraform compliant
 *
 * This service provider allows the creation of fake 'cdns' and 'lbs' resources
 *
 * API version: 1.0.0
 * Contact: apiteam@serviceprovider.io
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package api

import (
	"net/http"
)

func LBCreateV1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func LBDeleteV1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func LBGetV1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func LBUpdateV1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
