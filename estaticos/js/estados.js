/**
 * ESTADOS.JS - Lógica para Visualización y Reportes
 */

// ==========================================
// 1. UTILIDADES GLOBALES (Definir primero)
// ==========================================
window.mostrarAviso = function(mensaje) {
  const toastEl = document.getElementById('toastAtajo');
  const toastMsg = document.getElementById('toastMensaje');

  if (!toastEl || !toastMsg) {
    console.log("Aviso (consola):", mensaje);
    return;
  }

  toastMsg.innerHTML = `<i class="bi bi-info-circle me-2"></i> ${mensaje}`;
  const toast = bootstrap.Toast.getOrCreateInstance(toastEl);
  toast.show();
};

// ==========================================
// 2. FUNCIÓN GENERAR PDF (Ahora ya conoce mostrarAviso)
// =========================================

// Variable global para guardar el logo convertido
let logoDataURL = null;
const prepararLogo = () => {
  const svgString = `<?xml version="1.0" encoding="UTF-8"?><svg version="1.1" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024"><g transform="translate(0,1024) scale(0.1,-0.1)" fill="#28a745" stroke="none"><path d="M4805 9284 c-16-2-93-11-170-19-396-44-845-163-1152-305-24-11-46-20-49-20-11 0-267-131-371-191-169-96-428-274-558-383-305-258-483-440-702-721-210-269-438-662-544-940-11-27-23-59-28-70-71-167-157-460-196-665-9-47-20-107-26-135-59-303-77-800-39-1135 24-219 41-317 86-520 29-127 105-386 138-470 8-19 29-73 46-120 74-194 236-503 393-748 88-135 262-360 387-497 61-68 277-283 320-320 19-17 64-55 100-86 290-250 729-521 1060-654 36-14 74-30 85-35 114-54 389-139 630-195 495-114 1161-127 1655-30 558 108 1061 309 1510 603 396 260 724 560 1013 926 48 61 91 116 96 122 25 31 131 190 186 278 133 217 237 424 324 651 24 61 47 118 52 127 5 10 9 25 9 33 0 8 4 23 10 33 12 22 50 148 89 294 206 783 165 1680-110 2428-52 141-54 146-88 222-12 26-21 50-21 53 0 3-26 58-57 122-281 581-712 1109-1231 1511-107 82-346 242-467 313-148 86-373 196-510 249-44 17-93 37-110 44-16 8-48 19-70 26-22 7-49 16-60 20-31 12-148 48-170 53-11 3-63 16-115 30-185 50-454 97-674 117-100 9-600 12-671 4zm970-1254 c358-358 650-657 650-665 0-17-153-185-169-185-6 0-259 248-563 552-303 303-558 555-565 559-10 5-176-154-551-529-295-295-540-540-543-544-4-4-1-14 6-23 7-8 9-15 4-16-35-4-270 4-275 9-4 4-1 57 7 117 7 61 13 118 14 128 0 24 27 22 35-3 3-11 10-20 14-20 4 0 274 268 600 596 445 445 591 598 582 607-14 14-14 27 0 27 5 0 7-4 4-10-14-22 27-8 57 20 17 17 34 30 37 30 3 0 299-293 656-650zm-41-697 c332-335 608-612 612-616 9-8-164-187-181-187-5 0-242 232-525 515-283 283-518 514-522 513-4-2-341-336-748-743-598-597-736-740-718-742 13-1 90-4 171-5 l148-3 571 573 c314 314 573 572 577 572 4 0 240-234 525-519 l519-519 178 179 c98 98 182 179 187 179 12 0 1393-1388 1400-1408 1-4-631-639-1405-1412 l-1408-1405-609 611-609 612 89 91 88 92 523-523 523-523 745 745 c410 410 745 751 745 757 0 10-41 13-173 13 l-172 0-570-570 c-313-313-572-570-575-570-3 0-239 234-525 520 l-520 520-180-187-180-186-700 702 c-385 386-701 706-703 710-1 4 584 595 1300 1315 717 719 1346 1353 1398 1407 52 54 101 99 107 99 7 0 285-273 617-607zm1099-380 c68-6 81-13 82-45 0-9 394-413 875-897 481-485 876-885 878-889 1-4-402-411-897-906 l-898-898-82 92 c-44 51-81 95-81 99 0 3 361 367 803 809 l802 802-790 788 c-525 523-796 787-810 787-18 0-20 7-23 65-2 36 0 98 3 139 l7 73 31-6 c18-4 63-10 100-13zm-3373-110 c40-48 74-93 77-100 3-8-328-345-803-818 l-809-805 790-790 c596-596 794-788 807-784 17 5 18-6 18-135 0-97-4-141-11-141-6 0-41 7-78 15-36 8-83 17-103 20-40 6-46 13-26 33 9 9-183 207-867 892-484 484-880 885-880 890 0 16 1785 1810 1800 1810 8 0 46-39 85-87zm1093-4350 l568-568 545 547 c448 450 544 551 539 568-6 19-2 20 131 20 l137 0-7-42 c-3-24-9-80-13-125-3-46-9-83-12-83-3 0-12 10-20 21-12 21-43-8-657-622 l-644-644-650 650 c-358 357-650 654-650 659 0 9 23 38 104 129 28 31 53 57 56 57 3 0 261-255 573-567z"/><path d="M3160 5484 l0-236-89-91 c-49-51-99-104-111-119 l-21-28 110 0 111 0 0-249 c0-192 3-250 13-254 19-7 850-9 870-1 15 5 17 21 17 115 l0 109-304 0 c-234 0-305 3-308 13-2 6-3 68-3 137 l0 125 263 3 262 2 0 110 0 110-262 2-263 3 0 130 0 130 298 3 297 2 0 110 0 110-440 0-440 0 0-236z"/><path d="M4280 5480 c0-193-3-242-14-246-7-3-57-54-110-114 l-97-110 110 0 111 0 2-252 3-253 335-1 c184 0 386 2 448 4 l112 3 0 87 c0 48-3 97-6 110 l-6 22-304 0-304 0 0 140 0 140 265 0 265 0 0 110 0 110-265 0-266 0 3 133 3 132 250 1 c138 1 271 2 298 3 l47 1 0 110 0 110-440 0-440 0 0-240z"/><path d="M5372 5113 l3-608 135 0 135 0 3 242 2 243 263 2 262 3 3 108 3 107-193 1 c-106 1-224 2-263 3 l-70 1 0 140 0 140 299 3 299 2-7 53 c-3 28-6 78-6 110 l0 57-435 0-435 0 2-607z"/><path d="M6420 5110 l0-611 138 3 137 3 3 243 2 242 266 0 265 0-3 111-3 111-250-1 c-137 0-256 3-263 8-9 6-12 41-10 142 l3 134 298 3 297 2 0 110 0 110-440 0-440 0 0-610z"/></g></svg>`;
  const blob = new Blob([svgString], { type: 'image/svg+xml;charset=utf-8' });
  const url = URL.createObjectURL(blob);
  const img = new Image();
  img.onload = () => {
    const canvas = document.createElement('canvas');
    canvas.width = 1024; canvas.height = 1024;
    const ctx = canvas.getContext('2d');
    ctx.drawImage(img, 0, 0);
    logoDataURL = canvas.toDataURL('image/png');
    URL.revokeObjectURL(url);
  };
  img.src = url;
};
prepararLogo();

// 2. Función Principal
window.generarReportePDF = function() {
  try {
    window.mostrarAviso("Generando Reporte Ejecutivo...");
    const { jsPDF } = window.jspdf;

    // --- CONFIGURACIÓN INICIAL (PORTADA VERTICAL) ---
    const doc = new jsPDF('p', 'pt', 'letter');
    const fecha = new Date().toLocaleDateString('es-GT');
    const ancho = doc.internal.pageSize.width;
    const alto = doc.internal.pageSize.height;
    const colores = { primario: [40, 167, 69], oscuro: [33, 37, 41], info: [23, 162, 184] };

    // --- DISEÑO DE LA PORTADA ---
    doc.setFillColor(colores.primario[0], colores.primario[1], colores.primario[2]);
    doc.rect(0, 0, 50, alto, 'F');

    if (logoDataURL) {
      doc.addImage(logoDataURL, 'PNG', (ancho / 2) - 50, 150, 100, 100);
    }

    doc.setTextColor(colores.oscuro[0], colores.oscuro[1], colores.oscuro[2]);
    doc.setFontSize(28);
    doc.setFont('helvetica', 'bold');
    doc.text("ESTADOS FINANCIEROS", ancho / 2, 320, { align: 'center' });

    doc.setFontSize(18);
    doc.setFont('helvetica', 'normal');
    doc.text("Reporte Corporativo Anual", ancho / 2, 355, { align: 'center' });

    doc.setDrawColor(colores.primario[0], colores.primario[1], colores.primario[2]);
    doc.setLineWidth(2);
    doc.line((ancho / 2) - 100, 380, (ancho / 2) + 100, 380);

    doc.setFontSize(12);
    doc.text(`Fecha de Emisión: ${fecha}`, ancho / 2, 420, { align: 'center' });
    doc.text("Sistema Contable v0.01", ancho / 2, 440, { align: 'center' });

    doc.setFontSize(10);
    doc.setTextColor(100);
    doc.text("Este documento contiene información financiera confidencial.", ancho / 2, alto - 60, { align: 'center' });

    // --- FUNCIÓN HELPER PARA ENCABEZADOS POSTERIORES ---
    const agregarHeader = (titulo) => {
      const pageAncho = doc.internal.pageSize.width;
      doc.setFillColor(colores.oscuro[0], colores.oscuro[1], colores.oscuro[2]);
      doc.rect(0, 0, pageAncho, 50, 'F');
      if (logoDataURL) doc.addImage(logoDataURL, 'PNG', 20, 10, 30, 30);

      doc.setFontSize(16);
      doc.setTextColor(255, 255, 255);
      doc.setFont('helvetica', 'bold');
      doc.text(titulo, 60, 32);

      doc.setFontSize(9);
      doc.setFont('helvetica', 'normal');
      doc.text(`Generado: ${fecha}`, pageAncho - 120, 30);
    };

    // --- PÁGINA 2: HOJA DE TRABAJO (HORIZONTAL) ---
    doc.addPage('letter', 'l');
    agregarHeader("HOJA DE TRABAJO CONTABLE");
    doc.autoTable({
      html: "#hojaTrabajo",
      startY: 65,
      theme: 'grid',
      styles: { fontSize: 6.5, cellPadding: 2 },
      headStyles: { fillColor: colores.primario },
      didParseCell: function(data) {
        const rawElement = data.cell.raw;
        if (rawElement && (rawElement.classList.contains('border-top') || rawElement.classList.contains('ClaTotal'))) {
          data.cell.styles.lineWidth = { top: 1.5, bottom: 0, left: 0, right: 0 };
          data.cell.styles.fontStyle = 'bold';
        }
        if (rawElement && rawElement.classList.contains('border-double')) {
          data.cell.styles.lineWidth = { top: 0.5, bottom: 2, left: 0, right: 0 };
        }
      }
    });

    // --- PÁGINA 3: RESULTADOS (VERTICAL) ---
    doc.addPage('letter', 'p');
    agregarHeader("ESTADO DE RESULTADOS");
    doc.autoTable({
      html: "#resultados",
      startY: 65,
      theme: 'striped',
      headStyles: { fillColor: colores.primario },
      didParseCell: function(data) {
        const raw = data.cell.raw;
        if (raw && (raw.id === 'baseContable' || raw.classList.contains('fw-bold'))) {
          data.cell.styles.fontStyle = 'bold';
          data.cell.styles.fillColor = [240, 240, 240];
          data.cell.styles.lineWidth = { top: 1 };
        }
      }
    });

    // --- PÁGINA 4: BALANCE (VERTICAL) ---
    doc.addPage('letter', 'p');
    agregarHeader("BALANCE GENERAL");
    doc.autoTable({
      html: "#balance",
      startY: 65,
      theme: 'grid',
      headStyles: { fillColor: colores.primario },
      styles: { fontSize: 9 },
      didParseCell: function(data) {
        const raw = data.cell.raw;
        if (raw && (raw.classList.contains('fw-bold') || raw.classList.contains('border-top'))) {
          data.cell.styles.fontStyle = 'bold';
        }
      }
    });

    // Firmas
    let finalY = doc.lastAutoTable.finalY + 60;
    doc.autoTable({
      startY: finalY,
      theme: 'plain',
      body: [
        ['f.__________________________', 'f.__________________________'],
        ['Representante Legal', 'Contador General']
      ],
      styles: { halign: 'center', fontSize: 10 }
    });


    // --- PÁGINA 4: ANEXOS FISCALES E INDICADORES (TODO JUNTO) ---
    doc.addPage('letter', 'p');
    agregarHeader("ANEXOS Y ANÁLISIS");

    // 1. Tabla Fiscal
    doc.autoTable({
      startY: 70,
      margin: { left: 40 },
      tableWidth: 250,
      head: [['Descripción Fiscal', 'Monto']],
      body: [
        ['(+) ISR:', document.getElementById('txtISR')?.innerText || "0.00"],
        ['(+) Reserva Legal:', document.getElementById('txtReserva')?.innerText || "0.00"],
        ['(=) Utilidad Neta:', document.getElementById('txtNeto')?.innerText || "0.00"]
      ],
      headStyles: { fillColor: colores.info }
    });

    // 2. Tabla Indicadores (SIN doc.addPage para que use la misma hoja)
    const dataIndices = [];
    document.querySelectorAll('.card').forEach(card => {
      const n = card.querySelector('.card-subtitle')?.innerText.trim();
      const v = card.querySelector('.card-title')?.innerText.trim();
      const i = card.querySelector('.card-text')?.innerText.trim();
      if (n && v) dataIndices.push([n, v, i]);
    });

    doc.autoTable({
      // startY usa el final de la tabla anterior + 30 de espacio
      startY: doc.lastAutoTable.finalY + 30,
      head: [['Indicador', 'Valor', 'Interpretación']],
      body: dataIndices,
      headStyles: { fillColor: colores.oscuro },
      styles: { fontSize: 9 }
    });

    // --- NUMERACIÓN DE PÁGINAS (Insertar antes de doc.save) ---
    const totalPaginas = doc.internal.getNumberOfPages();
    for (let i = 1; i <= totalPaginas; i++) {
      doc.setPage(i);
      const pageAncho = doc.internal.pageSize.width;
      const pageAlto = doc.internal.pageSize.height;
      doc.setFontSize(9);
      doc.setTextColor(150);
      doc.setFont('helvetica', 'italic');
      // Texto: "Página X de Y"
      const textoPagina = `Página ${i} de ${totalPaginas}`;
      // Lo posicionamos en la esquina inferior derecha
      doc.text(textoPagina, pageAncho - 40, pageAlto - 20, { align: 'right' });
      // Opcional: Una pequeña línea divisoria gris en el pie de página
      doc.setDrawColor(200);
      doc.setLineWidth(0.5);
      doc.line(40, pageAlto - 30, pageAncho - 40, pageAlto - 30);
    }
    // --- GUARDADO FINAL ---
    doc.save(`Reporte_Contable.pdf`);
    window.mostrarAviso("¡PDF generado con éxito!");

  } catch (error) {
    console.error("ERROR GENERANDO PDF:", error);
    alert("Error al generar PDF: " + error.message);
  }
};

// ==========================================
// 1. CÁLCULOS FISCALES (Simulador)
// ==========================================
window.calcularFiscal = function(event) {
  if (event) event.preventDefault();
  console.log("Calculando ISR y Reserva...");

  const baseElement = document.getElementById('baseContable');
  if (!baseElement) return;

  const utilidadContable = parseFloat(baseElement.innerText.replace(/,/g, '').trim()) || 0;
  const porceISR = (parseFloat(document.getElementById('inputISR')?.value) || 0) / 100;
  const porceReserva = (parseFloat(document.getElementById('inputReserva')?.value) || 0) / 100;

  const isr = utilidadContable > 0 ? utilidadContable * porceISR : 0;
  const utilidadPreReserva = utilidadContable - isr;
  const reservaLegal = utilidadPreReserva > 0 ? utilidadPreReserva * porceReserva : 0;
  const netoFinal = utilidadContable - isr - reservaLegal;

  const f = new Intl.NumberFormat('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 });

  if (document.getElementById('txtISR')) document.getElementById('txtISR').innerText = f.format(isr);
  if (document.getElementById('txtReserva')) document.getElementById('txtReserva').innerText = f.format(reservaLegal);
  if (document.getElementById('txtNeto')) document.getElementById('txtNeto').innerText = f.format(netoFinal);

  window.mostrarAviso("Cálculos fiscales actualizados");
};

// ==========================================
// 2. MODAL DE AUDITORÍA
// ==========================================
window.verAuditoria = function(nombre, detalleCuenta, interpretacion) {
  const modalElemento = document.getElementById('modalAuditoria');
  const txtTitulo = document.getElementById('tituloAuditoria');
  const txtDesc = document.getElementById('descAuditoria');
  const txtFormula = document.getElementById('formulaAuditoria');

  if (!modalElemento) return;

  if (txtTitulo) txtTitulo.innerText = "Auditoría: " + nombre;
  if (txtDesc) txtDesc.innerText = interpretacion;
  if (txtFormula) txtFormula.innerHTML = `<span class="text-primary">${detalleCuenta}</span>`;

  const modalBS = bootstrap.Modal.getOrCreateInstance(modalElemento);
  modalBS.show();
};

