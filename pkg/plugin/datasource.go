package plugin

import (
	"context"
	"encoding/json"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/sqlds/v2"
	"net/http"
)

func NewDatasource(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	ds := sqlds.NewDatasource(&Ydb{})
	ds.CustomRoutes = setupCustomRoutes(settings)

	return ds.NewDatasource(settings)
}

func setupCustomRoutes(settings backend.DataSourceInstanceSettings) map[string]func(http.ResponseWriter, *http.Request) {
	return map[string]func(http.ResponseWriter, *http.Request){
		"/listTables": func(w http.ResponseWriter, r *http.Request) {
			if err := func(w http.ResponseWriter) error {
				tablesString, err := RetrieveListTablesForRoot(r.Context(), settings)
				if err != nil {
					return err
				}
				_, err = w.Write(tablesString)
				if err != nil {
					return err
				}
				return nil
			}(w); err != nil {
				w.WriteHeader(500)
				jsonErr, _ := json.Marshal(err.Error())
				w.Write(jsonErr)
				log.DefaultLogger.Error(err.Error())
			}
		},
		"/listFields": func(w http.ResponseWriter, r *http.Request) {
			if err := func(w http.ResponseWriter) error {
				query := r.URL.Query()
				table := query.Get("table")
				fieldsString, err := RetrieveTableFields(r.Context(), settings, table)
				if err != nil {
					return err
				}
				_, err = w.Write(fieldsString)
				if err != nil {
					return err
				}
				return nil
			}(w); err != nil {
				w.WriteHeader(500)
				jsonErr, _ := json.Marshal(err.Error())
				w.Write(jsonErr)
				log.DefaultLogger.Error(err.Error())
			}
		},
	}
}
