// ==========================================
// 1. VARIABLES GLOBALES Y CONFIGURACIÓN
// ==========================================
let currentInput = '0';
let subtotal = 0;
let pendingOp = null;
let ultimoInputEnfocado = null;
let yaSonó = false;

const sonidoExito = new Audio('https://assets.mixkit.co/active_storage/sfx/2568/2568-preview.mp3');
sonidoExito.volume = 0.3;

// ==========================================
// 2. FUNCIONES DE CÁLCULO Y FORMATEO (GLOBALES)
// ==========================================
const formatearVisual = (input) => {
  let val = parseFloat(input.value.replace(/,/g, ''));
  if (!isNaN(val)) {
    input.value = val.toLocaleString('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    });
  }
};

const calcularDiferencia = () => {
  let totalDebe = 0;
  let totalHaber = 0;
  let llenados = 0;
  const inputs = document.querySelectorAll('.input-contable');
  const diffValue = document.getElementById('diffValue');
  const cuadro = document.getElementById('cuadroCuadre');

  inputs.forEach(input => {
    const nameInput = input.getAttribute('name') || "";
    const valorLimpio = input.value.replace(/,/g, '');
    const val = parseFloat(valorLimpio) || 0;

    if (nameInput === "220001" || nameInput === "220005") return;

    if (val !== 0) {
      llenados++;
      const saldo = input.getAttribute('data-saldo');
      if (['activo', 'perdida', 'costo_debe'].includes(saldo)) totalDebe += val;
      else if (['pasivo', 'ganancia', 'costo_haber'].includes(saldo)) totalHaber += val;
    }
  });

  const diff = Math.abs(totalDebe - totalHaber);

  if (diffValue && cuadro) {
    diffValue.innerText = diff.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
    diffValue.style.fontSize = (diff > 999999) ? "1.2rem" : "1.5rem";

    if (diff < 0.01 && llenados > 0) {
      diffValue.style.color = "#00ff00";
      diffValue.innerHTML = '<i class="fas fa-check-circle"></i> CUADRADO';
      cuadro.style.borderColor = "#00ff00";
      cuadro.style.boxShadow = "0 0 20px rgba(0, 255, 0, 0.5)";
      if (!yaSonó) {
        sonidoExito.play().catch(() => { });
        yaSonó = true;
      }
    } else {
      diffValue.style.color = "#ff4444";
      cuadro.style.borderColor = "#ff4444";
      cuadro.style.boxShadow = "0 4px 15px rgba(0,0,0,0.5)";
      yaSonó = false;
    }
  }
};

// ==========================================
// 3. LÓGICA DE LA CALCULADORA
// ==========================================
function toggleCalc() {
  const calc = document.getElementById('calculadoraContable');
  const btnFlotante = document.getElementById('btnSumadora');
  if (!calc) return;

  if (calc.classList.contains('d-none')) {
    calc.classList.remove('d-none');
    if (btnFlotante) btnFlotante.style.display = 'none';
    document.querySelectorAll('.input-contable').forEach(i => {
      if (i.id !== 'input-220001') i.style.backgroundColor = "#1e1e1e";
    });
  } else {
    calc.classList.add('d-none');
    if (btnFlotante) btnFlotante.style.display = 'block';
  }
}

function updateDisplay() {
  const display = document.getElementById('displayCalc');
  if (display) display.value = currentInput;
}

function addNum(num) {
  if (currentInput === '0') currentInput = num;
  else currentInput += num;
  updateDisplay();
}

function ejecutarOperacion(nuevoValor) {
  if (pendingOp === '+') subtotal += nuevoValor;
  else if (pendingOp === '-') subtotal -= nuevoValor;
  else if (pendingOp === '*') subtotal *= nuevoValor;
  else if (pendingOp === '/') subtotal = nuevoValor !== 0 ? subtotal / nuevoValor : 0;
}

function setOp(op) {
  const valorActual = parseFloat(currentInput);
  if (pendingOp !== null) ejecutarOperacion(valorActual);
  else subtotal = valorActual;
  pendingOp = op;
  addToCinta(currentInput + ' ' + op);
  currentInput = '0';
  updateDisplay();
}

function calcular() {
  if (pendingOp === null) return;
  const valorFinal = parseFloat(currentInput);
  ejecutarOperacion(valorFinal);
  addToCinta(currentInput);
  addToCinta("----------");
  addToCinta("* " + subtotal.toLocaleString('en-US', { minimumFractionDigits: 2 }));
  currentInput = subtotal.toString();
  pendingOp = null;
  updateDisplay();
}

function pegarResultado() {
  const display = document.getElementById('displayCalc');
  if (display && ultimoInputEnfocado) {
    ultimoInputEnfocado.value = display.value;
    ultimoInputEnfocado.dispatchEvent(new Event('input'));
    addToCinta("-> Trasladado");
    ultimoInputEnfocado.focus();
  } else {
    mostrarAviso("Primero selecciona una celda");
  }
}

function addToCinta(text) {
  const cinta = document.getElementById('cintaContable');
  if (!cinta) return;
  const line = document.createElement('div');
  line.style.color = text.includes('-') ? "#ff4444" : "inherit";
  line.style.fontSize = "0.75rem";
  line.innerText = text;
  cinta.appendChild(line);
  cinta.scrollTop = cinta.scrollHeight;
}

function mostrarAviso(mensaje) {
  const toastEl = document.getElementById('toastAtajo');
  if (!toastEl) { console.log(mensaje); return; }
  document.getElementById('toastMensaje').innerHTML = `<i class="fas fa-keyboard me-2"></i> ${mensaje}`;
  new bootstrap.Toast(toastEl, { delay: 2000 }).show();
}

// ==========================================
// 4. INICIALIZACIÓN Y EVENTOS
// ==========================================
document.addEventListener('DOMContentLoaded', () => {
  const cuadro = document.getElementById('cuadroCuadre');
  const form = document.getElementById('balanceForm');
  document.getElementById('inputISR')?.addEventListener('input', calcularFiscal);
  document.getElementById('inputReserva')?.addEventListener('input', calcularFiscal);
  if (cuadro) document.body.appendChild(cuadro);

  // Sincronización de cuentas automáticas
  document.addEventListener('input', (e) => {
    if (e.target.classList.contains('input-contable')) {
      calcularDiferencia();
      // Lógica de espejos (111301 -> 220001, etc)
      const mappings = { '111301': '220001', '111303': '311001', '111304': '311002', '111305': '311001' };
      if (mappings[e.target.name]) {
        const target = document.getElementsByName(mappings[e.target.name])[0];
        if (target) { target.value = e.target.value; target.dispatchEvent(new Event('input')); }
      }
    }
  });

  document.addEventListener('focusin', (e) => {
    if (e.target.classList.contains('input-contable')) {
      ultimoInputEnfocado = e.target;
      e.target.value = e.target.value.replace(/,/g, '');
      e.target.select();
    }
  });

  document.addEventListener('focusout', (e) => {
    if (e.target.classList.contains('input-contable')) formatearVisual(e.target);
  });

  if (form) {
    form.addEventListener('submit', () => {
      document.querySelectorAll('.input-contable').forEach(i => i.value = i.value.replace(/,/g, ''));
    });
  }
  // Resaltar campos automáticos para no confundirlos con manuales
  const automaticos = ['220005', '311002'];
  automaticos.forEach(name => {
    const el = document.getElementsByName(name)[0];
    if (el) {
      el.style.borderLeft = "4px solid #00ff00";
      el.title = "Calculado automáticamente por el módulo fiscal";
    }
  });
  calcularDiferencia();
});

// Atajos de teclado (High Priority)
// ==========================================
// 2. ATAJOS DE TECLADO (VERSIÓN BLINDADA)
// ==========================================
document.addEventListener('keydown', (e) => {
  // 1. Validación de seguridad: Si no hay tecla, ignorar
  if (!e || !e.key) return;

  const key = e.key.toLowerCase();
  const isCtrl = e.ctrlKey || e.metaKey;

  // --- A. EL ATAJO MAESTRO: CTRL + S ---
  if (isCtrl && key === 's') {
    e.preventDefault();
    e.stopPropagation(); // Evita que otros scripts interfieran

    console.log("💾 Guardado rápido activado...");

    // Limpiar comas antes de enviar
    document.querySelectorAll('.input-contable').forEach(input => {
      input.value = input.value.replace(/,/g, '');
    });

    const form = document.getElementById('balanceForm') || document.querySelector('form');
    const diffValue = document.getElementById('diffValue');
    const estaCuadrado = diffValue && diffValue.innerText.includes("CUADRADO");

    if (form) {
      if (estaCuadrado) {
        mostrarAviso("¡Cuadrado! Guardando...");
        form.submit();
      } else {
        if (confirm("⚠️ La hoja NO está cuadrada. ¿Deseas guardar de todas formas?")) {
          mostrarAviso("Guardando con descuadre...");
          form.submit();
        }
      }
    }
    return;
  }

  // --- B. LÓGICA DE CALCULADORA ---
  const calc = document.getElementById('calculadoraContable');
  const estaVisible = calc && !calc.classList.contains('d-none');

  // Alt + C: Abrir/Cerrar
  if (e.altKey && key === 'c') {
    e.preventDefault();
    toggleCalc();
    return;
  }

  if (estaVisible) {
    if (/[0-9]/.test(key)) { e.preventDefault(); addNum(key); }
    if (['+', '-', '*', '/'].includes(key)) { e.preventDefault(); setOp(key); }
    if (key === '.' || key === ',') { e.preventDefault(); addDecimal(); }
    if (key === 'backspace') { e.preventDefault(); backspace(); }
    if (key === 'escape') { e.preventDefault(); toggleCalc(); }
    if (key === 'p') { e.preventDefault(); pegarResultado(); }
    if (key === 'enter') { e.preventDefault(); calcular(); }
  }
}, true);

// ==========================================
// FUNCIÓN FISCAL (Saca esto de cualquier otro bloque)
// ==========================================


window.calcularFiscal = function() {
  console.log("Calculando ISR y Reserva...");

  // 1. Obtener la utilidad contable (asegúrate de que este ID exista en tu HTML)
  const baseElement = document.getElementById('baseContable');
  if (!baseElement) {
    console.error("No se encontró el elemento 'baseContable'");
    return;
  }

  const utilidadContable = parseFloat(baseElement.innerText.replace(/,/g, '').trim()) || 0;

  // 2. Obtener porcentajes de los inputs (ajusta los IDs si son diferentes en tu HTML)
  const porceISR = (parseFloat(document.getElementById('inputISR')?.value) || 0) / 100;
  const porceReserva = (parseFloat(document.getElementById('inputReserva')?.value) || 0) / 100;

  // 3. Cálculos
  const isr = utilidadContable > 0 ? utilidadContable * porceISR : 0;
  const utilidadPreReserva = utilidadContable - isr;
  const reservaLegal = utilidadPreReserva > 0 ? utilidadPreReserva * porceReserva : 0;
  const netoFinal = utilidadContable - isr - reservaLegal;

  // 4. Formatear y mostrar resultados
  const f = new Intl.NumberFormat('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 });

  if (document.getElementById('txtISR')) document.getElementById('txtISR').innerText = f.format(isr);
  if (document.getElementById('txtReserva')) document.getElementById('txtReserva').innerText = f.format(reservaLegal);
  if (document.getElementById('txtNeto')) document.getElementById('txtNeto').innerText = f.format(netoFinal);

  mostrarAviso("Cálculos fiscales actualizados");
};

const btnPdf = document.getElementById("botonpdf");
if (btnPdf) {
  btnPdf.addEventListener("click", function() {
    mostrarAviso("Generando Reporte Profesional por Páginas...");
    const fecha = new Date().toLocaleDateString('es-GT');
    const nombreArchivo = `Estados_Financieros_${fecha.replace(/\//g, '-')}.pdf`;
    const { jsPDF } = window.jspdf;
    const doc = new jsPDF('l', 'pt', 'letter');
    const totalPaginasReales = 5;
    const agregarPiePagina = (doc, numero) => {
      const pageCount = doc.internal.getNumberOfPages();
      doc.setFontSize(9);
      doc.setTextColor(150);
      doc.text(`Página ${numero} de ${totalPaginasReales}`, doc.internal.pageSize.width - 80, doc.internal.pageSize.height - 30);
    };
    const configurarTabla = (idSelector, titulo, yInicia) => {
      doc.setFontSize(14);
      doc.setTextColor(0, 0, 0);
      doc.text(titulo, 40, yInicia);
      doc.autoTable({
        html: idSelector,
        startY: yInicia + 20,
        theme: 'grid',
        styles: { fontSize: 7.5, font: 'courier', cellPadding: 2, fillColor: [255, 255, 255], textColor: [0, 0, 0] },
        headStyles: { fillColor: [230, 230, 230], textColor: [0, 0, 0], fontStyle: 'bold' },
        didParseCell: function(data) {
          const el = data.cell.raw;
          if (el && el.tagName === 'INPUT') data.cell.text = el.value;
          if (el && el.classList?.contains('border-top')) data.cell.styles.lineWidth = { top: 1.5 };
          if (el && el.classList?.contains('border-double')) {
            data.cell.styles.lineWidth = { top: 1, bottom: 2.5 };
            data.cell.styles.fontStyle = 'bold';
          }
          if (data.column.index === 0) data.cell.styles.halign = 'left';
          else if (data.column.index > 1) data.cell.styles.halign = 'right';
        }
      });
    };
    configurarTabla("#hojaTrabajo", "HOJA DE TRABAJO - " + fecha, 50);
    agregarPiePagina(doc, 1)
    doc.addPage('letter', 'p');
    configurarTabla("#resultados", "ESTADO DE RESULTADOS", 50);
    agregarPiePagina(doc, 2)
    doc.addPage('letter', 'p'); // <--- SALTO DE PÁGINA PARA EL BALANCE
    configurarTabla("#balance", "BALANCE GENERAL", 50);
    doc.autoTable({
      startY: doc.lastAutoTable.finalY + 40,
      theme: 'plain',
      body: [['f.__________________________', 'f.__________________________'],
      ['Representante Legal', 'Contador General']],
      styles: { halign: 'center', fontSize: 9 }
    });
    agregarPiePagina(doc, 3)
    doc.addPage('letter', 'p');
    let finalY_BG = doc.lastAutoTable.finalY + 30;
    doc.setFontSize(11);
    doc.setFont('helvetica', 'bold');
    doc.text("ANEXO: CÁLCULO FISCAL", 40, finalY_BG);
    doc.autoTable({
      startY: finalY_BG + 10,
      margin: { left: 40 },
      tableWidth: 250,
      theme: 'plain',
      body: [
        ['(+) ISR:', document.getElementById('txtISR')?.innerText || "0.00"],
        ['(+) Reserva Legal:', document.getElementById('txtReserva')?.innerText || "0.00"],
        ['(=) Utilidad Neta:', document.getElementById('txtNeto')?.innerText || "0.00"]
      ],
      styles: { fontSize: 9, font: 'courier' }
    });
    agregarPiePagina(doc, 4)

    // --- PÁGINA 4: ANÁLISIS FINANCIERO (Vertical) ---
    doc.addPage('letter', 'p');
    doc.setFontSize(16);
    doc.setTextColor(0, 102, 204);
    doc.text("ANÁLISIS DE INDICADORES FINANCIEROS", 40, 50);
    // Diccionario de fórmulas para asegurar que aparezcan
    const diccionarioFormulas = {
      "LIQUIDEZ CORRIENTE": "Activo Corriente / Pasivo Corriente",
      "ROTACIÓN DE INVENTARIO": "Costo de Ventas / Promedio Inventarios",
      "DÍAS DE INVENTARIO": "365 / Rotación de Inventarios",
      "NIVEL DE ENDEUDAMIENTO": "(Total Pasivo / Total Activo) * 100",
      "MARGEN DE UTILIDAD": "(Utilidad Neta / Ventas Netas) * 100",
      "SOLVENCIA": "Total Activo / Total Pasivo"
    };
    const dataIndices = [];
    document.querySelectorAll('.card').forEach(card => {
      const nombreRaw = card.querySelector('.card-subtitle')?.innerText.trim() || "";
      const valor = card.querySelector('.card-title')?.innerText.trim() || "";
      const interpretacion = card.querySelector('.card-text')?.innerText.trim() || "";

      // Buscamos la fórmula en nuestro diccionario usando el nombre
      const formula = diccionarioFormulas[nombreRaw.toUpperCase()] || "Ver detalle en auditoría";

      if (nombreRaw && valor) {
        dataIndices.push([nombreRaw, valor, formula, interpretacion]);
      }
    });

    doc.autoTable({
      startY: 70,
      head: [['Indicador', 'Valor', 'Fórmula Aplicada', 'Interpretación']],
      body: dataIndices,
      theme: 'striped',
      headStyles: { fillColor: [44, 62, 80], textColor: [255, 255, 255] },
      columnStyles: {
        0: { cellWidth: 100, fontStyle: 'bold' },
        1: { cellWidth: 60, halign: 'center', fontStyle: 'bold' },
        2: { cellWidth: 120, fontSize: 8, font: 'courier' },
        3: { cellWidth: 'auto' }
      },
      styles: { fontSize: 8.5, cellPadding: 6 }
    });
    agregarPiePagina(doc, 5)

    doc.save(nombreArchivo);
    mostrarAviso("Reporte Maestro generado con éxito.");
  });
}
