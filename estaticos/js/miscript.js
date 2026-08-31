function activarSuscripcion(userId, e) {
  if (!confirm("¿Deseas activar la suscripción Pro para este usuario?")) return;

  // Feedback visual: deshabilitamos el botón para evitar doble clic
  const btn = e.currentTarget;
  const originalIcon = btn.innerHTML;
  btn.disabled = true;
  btn.innerHTML = '<i class="fas fa-spinner fa-spin"></i>';

  fetch(`/api/admin/activar-usuario/${userId}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
  })
    .then((response) => {
      if (!response.ok) throw new Error("Fallo en el servidor");
      // Usamos tu función mostrarAviso que ya tienes en el main js
      if (typeof mostrarAviso === "function") {
        mostrarAviso("¡Usuario activado exitosamente!");
      } else {
        alert("Usuario activado");
      }
      setTimeout(() => location.reload(), 1000);
    })
    .catch((error) => {
      alert("Error: " + error.message);
      btn.disabled = false;
      btn.innerHTML = originalIcon;
    });
}

function verUsuario(userId) {
  // Abrimos el modal (Bootstrap 5)
  const myModal = new bootstrap.Modal(document.getElementById("userModal"));
  myModal.show();

  // Pedimos los datos al servidor usando la ruta correcta del admin API
  fetch(`/api/admin/usuario/${userId}`)
    .then((res) => {
      if (!res.ok) throw new Error("Error al cargar datos");
      return res.json();
    })
    .then((data) => {
      const content = document.getElementById("modalContent");
      content.innerHTML = `
                <div class="list-group list-group-flush">
                    <div class="list-group-item"><strong>Nombre:</strong> ${data.Nombre} ${data.Apellido}</div>
                    <div class="list-group-item"><strong>Email:</strong> ${data.Email}</div>
                    <div class="list-group-item"><strong>Rol actual:</strong> <span class="badge bg-secondary">${data.Role}</span></div>
                    <div class="list-group-item"><strong>Prueba hasta:</strong> ${new Date(data.FechaFinPrueba).toLocaleDateString()}</div>
                    <div class="list-group-item"><strong>Suscripción:</strong>
                        ${data.SuscripcionActiva ? '<span class="text-success">Activa</span>' : '<span class="text-danger">Inactiva</span>'}
                    </div>
                </div>
            `;
    })
    .catch((err) => {
      document.getElementById("modalContent").innerHTML =
        '<p class="text-danger">Error al cargar datos</p>';
    });
}

function filtrarUsuarios() {
  const input = document.getElementById("busquedaUsuario");
  const filter = input.value.toLowerCase();
  const table = document.querySelector(".table");
  const tr = table.getElementsByTagName("tr");

  // Empezamos en 1 para saltarnos el encabezado (thead)
  for (let i = 1; i < tr.length; i++) {
    // Obtenemos el texto de la columna Nombre (índice 0) y Email (índice 1)
    const tdNombre = tr[i].getElementsByTagName("td")[0];
    const tdEmail = tr[i].getElementsByTagName("td")[1];

    if (tdNombre || tdEmail) {
      const txtNombre = tdNombre.textContent || tdNombre.innerText;
      const txtEmail = tdEmail.textContent || tdEmail.innerText;

      if (
        txtNombre.toLowerCase().indexOf(filter) > -1 ||
        txtEmail.toLowerCase().indexOf(filter) > -1
      ) {
        tr[i].style.display = ""; // Mostrar
      } else {
        tr[i].style.display = "none"; // Ocultar
      }
    }
    // Al final del bucle for
    let hayResultados = Array.from(tr)
      .slice(1)
      .some((row) => row.style.display !== "none");
    document.getElementById("noResultados").style.display = hayResultados
      ? "none"
      : "";
  }
}

function llenarModal(codigo, nombre, categoria, saldo) {
  document.getElementById("edit_codigo").value = codigo;
  document.getElementById("edit_nombre").value = nombre;
  document.getElementById("edit_categoria").value = categoria;
  document.getElementById("edit_saldo").value = saldo;
}

document.addEventListener("DOMContentLoaded", function () {
  const urlParams = new URLSearchParams(window.location.search);

  // Si la URL tiene ?success=actualizada
  if (urlParams.get("success") === "actualizada") {
    const toastLiveExample = document.getElementById("liveToast");
    const toast = new bootstrap.Toast(toastLiveExample);
    toast.show();

    // Limpiar la URL para que no vuelva a salir si refrescan la página
    window.history.replaceState({}, document.title, window.location.pathname);
  }
});

document.getElementById("buscarCuenta").addEventListener("keyup", function () {
  let val = this.value.toUpperCase();
  document.querySelectorAll("tbody tr").forEach((tr) => {
    tr.style.display = tr.innerText.toUpperCase().includes(val) ? "" : "none";
  });
});
