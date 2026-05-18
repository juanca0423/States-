import { useState } from 'react';
import { createRoot } from 'react-dom/client';

function FormularioContable() {
  const catalogoCompleto = window.CATALOGO_CUENTAS || [];

  const [cuentasSeleccionadas, setCuentasSeleccionadas] = useState([]);
  const [busqueda, setBusqueda] = useState('');

  // 1. Control de visibilidad del dropdown
  const [estaListaAbierta, setEstaListaAbierta] = useState(false);

  const agregarCuenta = (cuenta) => {
    const codigo = cuenta.Codigo || cuenta.codigo;
    if (!cuentasSeleccionadas.some(c => (c.Codigo || c.codigo) === codigo)) {
      setCuentasSeleccionadas([...cuentasSeleccionadas, { ...cuenta, Monto: '' }]);
    }

    // Al seleccionar, limpiamos y cerramos la lista
    setBusqueda('');
    setEstaListaAbierta(false);
  };

  const manejarMontoChange = (codigo, valor) => {
    setCuentasSeleccionadas(cuentasSeleccionadas.map(c => {
      const cCodigo = c.Codigo || c.codigo;
      return cCodigo === codigo ? { ...c, Monto: valor } : c;
    }));
  };

  const eliminarCuenta = (codigo) => {
    setCuentasSeleccionadas(cuentasSeleccionadas.filter(c => (c.Codigo || c.codigo) !== codigo));
  };

  // 2. Filtrado interactivo
  const cuentasFiltradas = estaListaAbierta ? catalogoCompleto.filter(c => {
    const nombreCuenta = c.Nombre || c.nombre || "";
    const codigoCuenta = c.Codigo || c.codigo || "";

    if (!busqueda) return true;

    return (
      nombreCuenta.toLowerCase().includes(busqueda.toLowerCase()) ||
      codigoCuenta.toString().includes(busqueda)
    );
  }) : [];

  return (
    <div>
      {/* SECTOR DE SELECT BUSCADOR INTELIGENTE (Se queda fijo arriba) */}
      <div className="mb-4 position-relative">
        <label className="text-success fw-bold mb-2">📁 Seleccionar Cuenta Contable:</label>

        <div className="input-group">
          <span className="input-group-text bg-dark border-success text-success">
            <i className="fas fa-list"></i>
          </span>
          <input
            type="text"
            className="form-control bg-dark text-white border-success"
            placeholder="Haz clic aquí o escribe para filtrar..."
            value={busqueda}
            onFocus={() => setEstaListaAbierta(true)}
            onChange={(e) => { setBusqueda(e.target.value); setEstaListaAbierta(true); }}
          />
          {busqueda && (
            <button
              className="btn btn-outline-secondary border-success text-muted"
              type="button"
              onClick={() => { setBusqueda(''); setEstaListaAbierta(false); }}
            >
              ❌
            </button>
          )}
        </div>

        {/* DESPLEGABLE TIPO SELECT */}
        <div
          className="list-group position-absolute w-100 shadow-lg overflow-auto"
          style={{
            zIndex: 1050,
            maxHeight: "250px",
            display: estaListaAbierta && cuentasFiltradas.length > 0 ? 'block' : 'none',
            border: "1px solid #28a745",
            borderRadius: "0.375rem"
          }}
        >
          {cuentasFiltradas.map(c => {
            const nombre = c.Nombre || c.nombre || "Sin Nombre";
            const codigo = c.Codigo || c.codigo || "000000";
            const saldo = c.Saldo || c.saldo || "";

            return (
              <button
                key={codigo}
                type="button"
                className="list-group-item list-group-item-action text-white d-flex justify-content-between align-items-center text-start py-2"
                style={{ fontSize: "0.9rem", backgroundColor: "#111", borderColor: "#222" }}
                onClick={() => agregarCuenta(c)}
              >
                <span>
                  <strong className="text-info">{codigo}</strong> — {nombre}
                </span>
                <span className="badge bg-success opacity-75">{saldo}</span>
              </button>
            );
          })}
        </div>
      </div>

      {/* ─── CONTENEDOR CON TAMAÑO FIJO REDUCIDO Y SCROLL ─── */}
      <div
        className="overflow-auto border border-secondary rounded p-2"
        style={{
          maxHeight: "380px",       // Reducido de 450px a 380px para ahorrar espacio
          backgroundColor: "#1a1d20"
        }}
      >
        <table className="table table-dark table-hover mb-0" style={{ borderCollapse: "separate" }}>
          <thead>
            <tr>
              {/* Forzamos el sticky directamente en las celdas <th> con fondo sólido */}
              <th
                style={{ width: "50%", position: "sticky", top: 0, backgroundColor: "#212529", zIndex: 10 }}
                className="text-success"
              >
                Cuenta Contable
              </th>
              <th
                style={{ width: "40%", position: "sticky", top: 0, backgroundColor: "#212529", zIndex: 10 }}
                className="text-success text-center"
              >
                Monto
              </th>
              <th
                style={{ width: "10%", position: "sticky", top: 0, backgroundColor: "#212529", zIndex: 10 }}
                className="text-success text-center"
              >
                Acción
              </th>
            </tr>
          </thead>
          <tbody>
            {cuentasSeleccionadas.length === 0 ? (
              <tr>
                <td colSpan="3" className="text-center text-muted py-4">
                  El formulario está vacío. Utiliza el buscador de arriba para añadir las cuentas del ejercicio.
                </td>
              </tr>
            ) : (
              cuentasSeleccionadas.map(c => {
                const nombre = c.Nombre || c.nombre || "Sin Nombre";
                const codigo = c.Codigo || c.codigo || "000000";
                const saldo = c.Saldo || c.saldo || "";
                const categoria = c.Categoria || c.categoria || "";

                return (
                  <tr key={codigo}>
                    <td className="align-middle text-start ps-3">
                      <span className="text-white fw-bold">{nombre}</span>
                      <small className="text-muted d-block">{codigo} {categoria && `(${categoria})`}</small>
                    </td>
                    <td>
                      <input
                        type="text"
                        name={codigo}
                        id={`input-${codigo}`}
                        data-saldo={saldo}
                        value={c.Monto}
                        onChange={(e) => manejarMontoChange(codigo, e.target.value)}
                        className="form-control monto-contable input-contable w-100 bg-transparent text-white border-info text-end"
                        placeholder="0.00"
                      />
                    </td>
                    <td className="text-center align-middle">
                      <button
                        type="button"
                        className="btn btn-outline-danger btn-sm"
                        onClick={() => eliminarCuenta(codigo)}
                      >
                        <i className="fas fa-trash-alt"></i>
                      </button>
                    </td>
                  </tr>
                );
              })
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}

// ─── EL INICIALIZADOR DEBE IR AQUÍ AFUERA (AL FINAL DEL ARCHIVO) ───
const container = document.getElementById('react-formulario-root');
if (container) {
  const root = createRoot(container);
  root.render(<FormularioContable />);
}
