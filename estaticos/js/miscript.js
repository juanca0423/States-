let balance = document.getElementById("balance");
balance.addEventListener("submit", function (e) {
  e.preventDefault();
  const form = e.target;
  const formData = new FormData(form);
  const data = {};
  formData.forEach((value, key) => {
    data[key] = value;
  });
  fetch("/estados", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error("Error en la respuesta del servidor");
      }
      return response.json();
    })
    .then((data) => {
      console.log("Éxito:", data);
      alert("¡Datos enviados con éxito!" + JSON.stringify(data));
    })
    .catch((error) => {
      console.error("Error:", error);
      alert("Error al enviar los datos." + JSOM.stringify(data));
    });
});
