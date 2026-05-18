// Funciones que se llaman desde el HTML (Globales)
function llenarModal(
  codigo,
  nombre,
  categoria,
  saldo,
  escosto,
  esvariable,
  esefectivo,
) {
  document.getElementById("edit_codigo").value = codigo;
  document.getElementById("edit_nombre").value = nombre;
  document.getElementById("edit_categoria").value = categoria;

  // Asignar Saldo
  const selectSaldo = document.getElementById("edit_saldo");
  if (selectSaldo) selectSaldo.value = saldo;

  // Asignar Booleans (convertidos a string para el select)
  document.querySelector('#modalEditar select[name="escosto"]').value =
    escosto.toString();
  document.getElementById("edit_es_variable").value = esvariable.toString();
  document.getElementById("edit_es_efectivo").value = esefectivo.toString();
}

document.addEventListener("DOMContentLoaded", () => {
  const inputBusqueda = document.getElementById("buscarCuenta");

  if (inputBusqueda) {
    inputBusqueda.addEventListener("keyup", function () {
      const val = this.value.toLowerCase(); // Buscamos en minúsculas para mayor coincidencia
      const filas = document.querySelectorAll("table tbody tr");

      filas.forEach((tr) => {
        // Obtenemos el texto de las celdas de Código (0) y Nombre (1)
        const codigo = tr.cells[0].textContent.toLowerCase();
        const nombre = tr.cells[1].textContent.toLowerCase();

        // El buscador ahora es insensible a mayúsculas/minúsculas
        if (codigo.includes(val) || nombre.includes(val)) {
          tr.style.display = "";
        } else {
          tr.style.display = "none";
        }
      });
    });
  }
});
