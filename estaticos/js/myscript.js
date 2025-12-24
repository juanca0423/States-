const real = document.getElementById("botonpdf");
real.addEventListener("click", function () {
  console.log("hola");
  const { jsPDF } = window.jspdf;
  const doc = new jsPDF();
  const tablas = ["#hojaTrabajo", "#resultados", "#balance"];
  tablas.forEach((idSelector, index) => {
    if (index > 0) doc.addPage();
    doc.autoTable({
      html: idSelector,
      theme: "striped", // Estilo visual
      headStyles: { fillColor: [41, 128, 185] }, // Color del encabezado
      margin: { top: 20 },
      didDrawPage: function (data) {
        // Opcional: Agregar un título manual en cada página
        doc.text("Estados Financieros", 14, 15);
      },
    });
  });
  doc.save("estados_financieros.pdf");
});
