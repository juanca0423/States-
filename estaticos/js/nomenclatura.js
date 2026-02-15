// Funciones que se llaman desde el HTML (Globales)
function llenarModal(codigo, nombre, categoria, saldo, escosto) {
  document.getElementById('edit_codigo').value = codigo;
  document.getElementById('edit_nombre').value = nombre;
  document.getElementById('edit_categoria').value = categoria;

  // IMPORTANTE: Asegúrate de tener este ID en el select del modal
  const selectSaldo = document.getElementById('edit_saldo');
  if (selectSaldo) {
    selectSaldo.value = saldo;
  }

  // Para el booleano
  const selectCosto = document.querySelector('#modalEditar select[name="escosto"]');
  if (selectCosto) {
    selectCosto.value = escosto.toString();
  }
}

document.addEventListener('DOMContentLoaded', () => {
  // Buscador rápido de la tabla
  const inputBusqueda = document.getElementById('buscarCuenta');
  if (inputBusqueda) {
    inputBusqueda.addEventListener('keyup', function() {
      let val = this.value.toUpperCase();

      const filas = document.querySelectorAll('table tbody tr');
      filas.forEach(tr => {
        const codigo = tr.cells[0].textContent.toUpperCase();
        const nombre = tr.cells[1].textContent.toUpperCase();
        tr.style.display = (codigo.incldes(val) || nombre.includes(val)) ? '' : 'none';
      });
    });
  }

  // Manejo del Toast de éxito (vía URL)
  const urlParams = new URLSearchParams(window.location.search);
  if (urlParams.get('success') === 'actualizada') {
    const toastEl = document.getElementById('liveToast');
    if (toastEl) new bootstrap.Toast(toastEl).show();
    window.history.replaceState({}, document.title, window.location.pathname);
  }
});
