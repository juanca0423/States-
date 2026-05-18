package config

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func ObtenerCuentasDesdeDB() ([]CuentaItem, error) {
	// La URL que copiaste de Supabase
	connStr := "postgresql://postgres:[#bEAqJBwne-7Ak7]@db.rmxmcnygzqranlcoinbr.supabase.co:5432/postgres"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT codigo, nombre, saldo FROM nomenclatura ORDER BY codigo ASC")
	if err != nil {
		return nil, err
	}

	var cuentas []CuentaItem
	for rows.Next() {
		var c CuentaItem
		// Escaneamos los datos de la DB a tu estructura
		err := rows.Scan(&c.Codigo, &c.Nombre, &c.Saldo)
		if err != nil {
			continue
		}
		cuentas = append(cuentas, c)
	}
	return cuentas, nil
}
