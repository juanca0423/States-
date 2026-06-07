// ==========================================
// 1. VARIABLES GLOBALES E INSTANCIAS
// ==========================================
let graficoEquilibrioInstance = null;
let graficoDonaInstance = null;
let graficoEstructuraInstance = null;
let graficoElementosCostoInstance = null; // NUEVA: Para la terna industrial
let logoDataURL = null;

// Para guardar las imágenes "en frío" para Excel
let imgEquilibrioBase64 = null;
let imgDonaBase64 = null;
let imgElementosCostoBase64 = null; // NUEVA

// Plugin para fondo blanco
const pluginFondoBlanco = {
  id: "customCanvasBackgroundColor",
  beforeDraw: (chart) => {
    const { ctx } = chart;
    ctx.save();
    ctx.globalCompositeOperation = "destination-over";
    ctx.fillStyle = "white";
    ctx.fillRect(0, 0, chart.width, chart.height);
    ctx.restore();
  },
};

// ==========================================
// 2. FUNCIÓN MAESTRA DE GRÁFICAS
// ==========================================
function renderizarGraficas() {
  const contenedor = document.getElementById("datos-financieros");
  if (!contenedor) return;

  // 1. Captura de datos REALES desde el servidor (vía dataset)
  const vNetas = parseFloat(contenedor.dataset.ventas) || 0;
  const cFijos = parseFloat(contenedor.dataset.fijos) || 0;
  const cVar = parseFloat(contenedor.dataset.variables) || 0;
  const pECaja = parseFloat(contenedor.dataset.puntoEfe) || 0;
  const pEContable = parseFloat(contenedor.dataset.puntoE) || 0;

  // Captura específica industrial (NUEVA)
  const mMP = parseFloat(contenedor.dataset.mp) || 0;
  const mMOD = parseFloat(contenedor.dataset.mod) || 0;
  const mCIF = parseFloat(contenedor.dataset.cif) || 0;

  // --- 1. GRÁFICA DE EQUILIBRIO (Lineal) ---
  const ctxEq = document.getElementById("graficoEquilibrio");
  if (ctxEq) {
    if (graficoEquilibrioInstance) graficoEquilibrioInstance.destroy();
    graficoEquilibrioInstance = new Chart(ctxEq.getContext("2d"), {
      type: "line",
      plugins: [pluginFondoBlanco],
      data: {
        labels: ["0%", "25%", "50%", "75%", "100%"],
        datasets: [
          {
            label: "Ingresos Totales",
            data: [0, vNetas * 0.25, vNetas * 0.5, vNetas * 0.75, vNetas],
            borderColor: "#FFD700",
            backgroundColor: "rgba(255, 215, 0, 0.1)",
            fill: true,
            tension: 0.1,
          },
          {
            label: "Costos Totales",
            data: [
              cFijos,
              cFijos + cVar * 0.25,
              cFijos + cVar * 0.5,
              cFijos + cVar * 0.75,
              cFijos + cVar,
            ],
            borderColor: "#FF6384",
            borderDash: [5, 5],
            fill: false,
            tension: 0.1,
          },
          {
            label: "Límite Equilibrio Caja",
            data: Array(5).fill(pECaja),
            borderColor: "#4BC0C0",
            borderWidth: 2,
            pointRadius: 0,
            fill: false,
            borderDash: [2, 2],
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        animation: { duration: 1000, easing: "easeInOutQuart" },
        scales: {
          y: { beginAtZero: true, ticks: { color: "#fff" } },
          x: { ticks: { color: "#fff" } },
        },
      },
    });
  }

  // --- 2. GRÁFICA DE DONA (Estructura de Costos) ---
  const ctxDo = document.getElementById("graficoDona");
  if (ctxDo) {
    if (graficoDonaInstance) graficoDonaInstance.destroy();
    graficoDonaInstance = new Chart(ctxDo.getContext("2d"), {
      type: "doughnut",
      data: {
        labels: ["Fijos", "Variables"],
        datasets: [
          {
            data: [cFijos, cVar],
            backgroundColor: ["#FF6384", "#36A2EB"],
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        animation: { animateRotate: true, animateScale: true },
      },
    });
  }

  // --- 3. GRÁFICA DE ESTRUCTURA (Comparativa de Barras) ---
  const ctxEst = document.getElementById("graficoEstructura");
  if (ctxEst) {
    if (graficoEstructuraInstance) graficoEstructuraInstance.destroy();
    graficoEstructuraInstance = new Chart(ctxEst.getContext("2d"), {
      type: "bar",
      data: {
        labels: ["Equilibrio Contable", "Equilibrio Caja", "Ventas Actuales"],
        datasets: [
          {
            label: "Monto Q",
            data: [pEContable, pECaja, vNetas],
            backgroundColor: ["#FF6384", "#4BC0C0", "#FFD700"],
            borderWidth: 1,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { display: false } },
        scales: {
          y: { beginAtZero: true, ticks: { color: "#fff" } },
          x: { ticks: { color: "#fff" } },
        },
      },
    });
  }

  // --- 4. NUEVA GRÁFICA: DISTRIBUCIÓN INDUSTRIAL (Pastel) ---
  const ctxCostosInd = document.getElementById("graficoElementosCosto");
  if (ctxCostosInd) {
    if (graficoElementosCostoInstance) graficoElementosCostoInstance.destroy();
    graficoElementosCostoInstance = new Chart(ctxCostosInd.getContext("2d"), {
      type: "pie",
      plugins: [pluginFondoBlanco],
      data: {
        labels: [
          "Materia Prima Consumida",
          "Mano de Obra Directa (MOD)",
          "Carga Fabril (CIF)",
        ],
        datasets: [
          {
            data: [mMP, mMOD, mCIF],
            backgroundColor: ["#36A2EB", "#FFCE56", "#9966FF"],
            borderWidth: 1,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            position: "bottom",
            labels: { color: "#fff" },
          },
        },
        animation: { animateRotate: true, animateScale: true },
      },
    });
  }
}

// ==========================================
// 3. UTILIDADES Y CARGA
// ==========================================
window.mostrarAviso = (mensaje) => console.log("🔔 Reporte:", mensaje);

function toggleCarga(idBoton, cargando) {
  const btn = document.getElementById(idBoton);
  if (!btn) return;
  const spinner =
    btn.querySelector(".spinner-border") ||
    document.getElementById("spinnerLogo") ||
    document.getElementById("spinnerExcel");
  if (cargando) {
    btn.disabled = true;
    if (spinner) spinner.classList.remove("d-none");
  } else {
    btn.disabled = false;
    if (spinner) spinner.classList.add("d-none");
  }
}

// ==========================================
// 4. CÁLCULOS FISCALES
// ==========================================
window.calcularFiscal = function (event) {
  if (event) event.preventDefault();
  const baseElement = document.getElementById("baseContable");
  if (!baseElement) return;

  const utilidadContable =
    parseFloat(baseElement.innerText.replace(/,/g, "").trim()) || 0;
  const porceISR =
    (parseFloat(document.getElementById("inputISR")?.value) || 0) / 100;
  const porceReserva =
    (parseFloat(document.getElementById("inputReserva")?.value) || 0) / 100;

  const isr = utilidadContable > 0 ? utilidadContable * porceISR : 0;
  const utilidadPreReserva = utilidadContable - isr;
  const reservaLegal =
    utilidadPreReserva > 0 ? utilidadPreReserva * porceReserva : 0;
  const netoFinal = utilidadContable - isr - reservaLegal;

  const f = new Intl.NumberFormat("en-US", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });
  if (document.getElementById("txtISR"))
    document.getElementById("txtISR").innerText = f.format(isr);
  if (document.getElementById("txtReserva"))
    document.getElementById("txtReserva").innerText = f.format(reservaLegal);

  const elNeto = document.getElementById("txtNeto");
  if (elNeto) {
    elNeto.innerText = f.format(netoFinal);
    elNeto.className =
      "fw-bold p-2 rounded " +
      (netoFinal > 0
        ? "ganancia-positiva"
        : netoFinal < 0
          ? "perdida-negativa"
          : "bg-secondary text-white");
  }
  window.mostrarAviso("Cálculos fiscales actualizados");
};

// ==========================================
// 5. EXPORTACIÓN Y LOGO
// ==========================================
const prepararLogo = async () => {
  try {
    const response = await fetch("/static/images/logoP.svg");
    const svgText = await response.text();
    const blob = new Blob([svgText], { type: "image/svg+xml" });
    const url = URL.createObjectURL(blob);
    const img = new Image();
    img.onload = () => {
      const canvas = document.createElement("canvas");
      canvas.width = 500;
      canvas.height = 500;
      const ctx = canvas.getContext("2d");
      ctx.drawImage(img, 0, 0, 500, 500);
      logoDataURL = canvas.toDataURL("image/png");
      URL.revokeObjectURL(url);
      const btn = document.getElementById("btnPDF");
      if (btn) btn.disabled = false;
    };
    img.src = url;
  } catch (err) {
    console.error("Error logo:", err);
  }
};

window.generarReportePDF = function () {
  try {
    toggleCarga("btnPDF", true);
    window.mostrarAviso("Generando Reporte Ejecutivo...");
    const { jsPDF } = window.jspdf;

    const doc = new jsPDF("p", "pt", "letter");
    const fecha = new Date().toLocaleDateString("es-GT");
    const ancho = doc.internal.pageSize.width;
    const alto = doc.internal.pageSize.height;
    const colores = { primario: [40, 167, 69], oscuro: [33, 37, 41] };

    const estilarCelda = (data) => {
      const raw = data.cell.raw;
      if (!raw) return;
      const texto = data.cell.text[0] || "";
      if (/^[\d,.-]+$/.test(texto.trim()) && texto.trim() !== "-") {
        data.cell.styles.halign = "right";
      }
      if (
        raw.classList.contains("border-top") ||
        raw.classList.contains("fw-bold")
      ) {
        data.cell.styles.lineWidth = { top: 1, bottom: 0, left: 0, right: 0 };
        data.cell.styles.fontStyle = "bold";
      }
      if (
        raw.classList.contains("border-bottom-double") ||
        raw.classList.contains("ClaTotal")
      ) {
        data.cell.styles.lineWidth = { top: 1, bottom: 2.5, left: 0, right: 0 };
        data.cell.styles.lineColor = [0, 0, 0];
        data.cell.styles.fontStyle = "bold";
      }
    };

    // --- 1. PORTADA ---
    doc.setFillColor(
      colores.primario[0],
      colores.primario[1],
      colores.primario[2],
    );
    doc.rect(0, 0, 50, alto, "F");
    if (logoDataURL)
      doc.addImage(logoDataURL, "PNG", ancho / 2 - 50, 150, 100, 100);
    doc.setFontSize(28);
    doc.setFont("helvetica", "bold");
    doc.text("ESTADOS FINANCIEROS", ancho / 2, 320, { align: "center" });

    const agregarHeader = (titulo) => {
      const pAncho = doc.internal.pageSize.width;
      doc.setFillColor(colores.oscuro[0], colores.oscuro[1], colores.oscuro[2]);
      doc.rect(0, 0, pAncho, 50, "F");
      doc.setFontSize(14);
      doc.setTextColor(255, 255, 255);
      doc.text(titulo, 60, 32);
      doc.setFontSize(9);
      doc.text(`Generado: ${fecha}`, pAncho - 120, 32);
    };

    // --- 2. TABLAS ---
    doc.addPage("letter", "l");
    agregarHeader("HOJA DE TRABAJO CONTABLE");
    doc.autoTable({
      html: "#hojaTrabajo",
      startY: 65,
      theme: "grid",
      styles: { fontSize: 6 },
      didParseCell: estilarCelda,
    });

    doc.addPage("letter", "p");
    agregarHeader("ESTADO DE RESULTADOS");
    doc.autoTable({
      html: "#resultados",
      startY: 65,
      theme: "striped",
      didParseCell: estilarCelda,
    });

    doc.addPage("letter", "p");
    agregarHeader("BALANCE GENERAL");
    doc.autoTable({
      html: "#balance",
      startY: 65,
      theme: "grid",
      didParseCell: estilarCelda,
    });

    // --- 3. FIRMAS ---
    let yF = doc.lastAutoTable.finalY + 40;
    if (yF > alto - 80) {
      doc.addPage();
      yF = 70;
    }
    doc.autoTable({
      startY: yF,
      theme: "plain",
      body: [
        ["f.__________________________", "f.__________________________"],
        ["Representante Legal", "Contador General"],
      ],
      styles: { halign: "center", fontStyle: "bold" },
    });

    // --- 4. ANEXOS E INDICADORES ---
    doc.addPage("letter", "p");
    agregarHeader("ANÁLISIS FINANCIERO");

    const dataIndices = [];
    document.querySelectorAll(".card").forEach((card) => {
      const n = card.querySelector(".card-subtitle")?.innerText.trim();
      const v = card.querySelector(".card-title")?.innerText.trim();
      if (n && v) dataIndices.push([n, v]);
    });

    doc.autoTable({
      startY: 70,
      head: [["Indicador", "Valor"]],
      body: dataIndices,
      headStyles: { fillColor: colores.oscuro },
    });

    // --- 5. GRÁFICOS (Soporte dinámico para el de fábrica) ---
    const canEq = document.getElementById("graficoEquilibrio");
    const canDo = document.getElementById("graficoDona");
    const canInd = document.getElementById("graficoElementosCosto");

    if (canEq || canDo || canInd) {
      doc.addPage("letter", "p");
      agregarHeader("DASHBOARD GRÁFICO");
      let yGr = 70;

      if (canEq) {
        doc.setFontSize(10);
        doc.setTextColor(0);
        doc.text("Punto de Equilibrio", 40, yGr);
        doc.addImage(
          canEq.toDataURL("image/png"),
          "PNG",
          40,
          yGr + 10,
          520,
          180,
        );
        yGr += 210;
      }

      if (canDo) {
        doc.text("Estructura de Gastos", 40, yGr);
        doc.addImage(
          canDo.toDataURL("image/png"),
          "PNG",
          40,
          yGr + 10,
          240,
          180,
        );
      }

      if (canInd) {
        doc.text("Elementos de Costo Industrial", ancho / 2 + 20, yGr);
        doc.addImage(
          canInd.toDataURL("image/png"),
          "PNG",
          ancho / 2 + 20,
          yGr + 10,
          240,
          180,
        );
      }
    }

    // --- 6. NUMERACIÓN FINAL ---
    const totalP = doc.internal.getNumberOfPages();
    for (let i = 1; i <= totalP; i++) {
      doc.setPage(i);
      doc.setFontSize(8);
      doc.setTextColor(150);
      doc.text(`Página ${i} de ${totalP}`, ancho - 40, alto - 20, {
        align: "right",
      });
      doc.setDrawColor(200);
      doc.line(40, alto - 30, ancho - 40, alto - 30);
    }

    doc.save(`Reporte_Final_${fecha}.pdf`);
    window.mostrarAviso("¡Reporte completo generado!");
  } catch (e) {
    console.error(e);
    alert("Error: " + e.message);
  } finally {
    toggleCarga("btnPDF", false);
  }
};

// ==========================================
// NUEVA FUNCIÓN DE ENLACE A EXCELJS
// ==========================================
window.descargarExcel = async function () {
  try {
    toggleCarga("btnExcel", true);
    window.mostrarAviso("Sincronizando datos con Excel...");

    if (!imgEquilibrioBase64) {
      imgEquilibrioBase64 = document
        .getElementById("graficoEquilibrio")
        ?.toDataURL("image/png");
    }
    if (!imgDonaBase64) {
      imgDonaBase64 = document
        .getElementById("graficoDona")
        ?.toDataURL("image/png");
    }
    // Salvaguarda para el gráfico industrial
    if (document.getElementById("graficoElementosCosto")) {
      imgElementosCostoBase64 = document
        .getElementById("graficoElementosCosto")
        .toDataURL("image/png");
      window.imagenIndustrialParaExcel = imgElementosCostoBase64;
    }

    window.imagenEquilibrioParaExcel = imgEquilibrioBase64;
    window.imagenDonaParaExcel = imgDonaBase64;

    if (typeof generarReporteExcelFinal === "function") {
      await generarReporteExcelFinal();
    }

    window.mostrarAviso("¡Excel generado con éxito!");
  } catch (error) {
    console.error("Error crítico en Excel:", error);
  } finally {
    toggleCarga("btnExcel", false);
  }
};

// ==========================================
// 6. INICIALIZACIÓN FINAL
// ==========================================
document.addEventListener("DOMContentLoaded", () => {
  const formateador = new Intl.NumberFormat("en-US", {
    minimumFractionDigits: 2,
  });
  document.querySelectorAll(".formatear-numero").forEach((celda) => {
    const valor = parseFloat(celda.textContent.trim());
    if (!isNaN(valor)) celda.textContent = formateador.format(valor);
  });

  document
    .getElementById("inputISR")
    ?.addEventListener("input", () => window.calcularFiscal());
  document
    .getElementById("inputReserva")
    ?.addEventListener("input", () => window.calcularFiscal());

  prepararLogo();
  renderizarGraficas();
  window.calcularFiscal();
});

window.verAuditoria = function (nombre, detalleCuenta, interpretacion) {
  const modalElemento = document.getElementById("modalAuditoria");
  if (!modalElemento) return;
  document.getElementById("tituloAuditoria").innerText = "Auditoría: " + nombre;
  document.getElementById("descAuditoria").innerText = interpretacion;
  document.getElementById("formulaAuditoria").innerText = detalleCuenta;
  bootstrap.Modal.getOrCreateInstance(modalElemento).show();
};
