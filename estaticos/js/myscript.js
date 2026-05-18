// ==========================================
// 1. VARIABLES GLOBALES (Siempre al inicio)
// ==========================================
let currentInput = "0";
let subtotal = 0;
let pendingOp = null;
let ultimoInputEnfocado = null;
let yaSonó = false;

const sonidoExito = new Audio(
  "https://assets.mixkit.co/active_storage/sfx/2568/2568-preview.mp3",
);
sonidoExito.volume = 0.3;

// ==========================================
// 2. EXPORTACIÓN DE FUNCIONES AL ÁMBITO GLOBAL (window)
// ==========================================

window.updateDisplay = function () {
  const display = document.getElementById("displayCalc");
  if (display) display.value = currentInput;
};

window.addNum = function (num) {
  if (currentInput === "0") currentInput = num;
  else currentInput += num;
  window.updateDisplay();
};

window.addDecimal = function () {
  if (!currentInput.includes(".")) {
    currentInput += currentInput === "" ? "0." : ".";
    window.updateDisplay();
  }
};

window.backspace = function () {
  if (currentInput.length > 1) {
    currentInput = currentInput.slice(0, -1);
  } else {
    currentInput = "0";
  }
  window.updateDisplay();
};

window.clearCalc = function () {
  currentInput = "0";
  subtotal = 0;
  pendingOp = null;
  const cinta = document.getElementById("cintaContable");
  if (cinta)
    cinta.innerHTML =
      '<div class="text-muted small text-center border-bottom mb-1">REGISTRO</div>';
  window.updateDisplay();
};

window["setOp"] = function (op) {
  const valorActual = parseFloat(currentInput);
  if (pendingOp !== null) window.ejecutarOperacion(valorActual);
  else subtotal = valorActual;
  pendingOp = op;
  window.addToCinta(currentInput + " " + op);
  currentInput = "0";
  window.updateDisplay();
};

window.ejecutarOperacion = function (nuevoValor) {
  if (pendingOp === "+") subtotal += nuevoValor;
  else if (pendingOp === "-") subtotal -= nuevoValor;
  else if (pendingOp === "*") subtotal *= nuevoValor;
  else if (pendingOp === "/")
    subtotal = nuevoValor !== 0 ? subtotal / nuevoValor : 0;
};

window.calcular = function () {
  if (pendingOp === null) return;
  const valorFinal = parseFloat(currentInput);
  window.ejecutarOperacion(valorFinal);
  window.addToCinta(currentInput);
  window.addToCinta("----------");
  window.addToCinta(
    "* " + subtotal.toLocaleString("en-US", { minimumFractionDigits: 2 }),
  );
  currentInput = subtotal.toString();
  pendingOp = null;
  window.updateDisplay();
};

window.toggleCalc = function () {
  const calc = document.getElementById("calculadoraContable");
  const btnFlotante = document.getElementById("btnSumadora");
  if (!calc) return;

  if (calc.classList.contains("d-none")) {
    calc.classList.remove("d-none");
    if (btnFlotante) btnFlotante.style.style.display = "none";
    document.querySelectorAll(".input-contable").forEach((i) => {
      if (i.id !== "input-220001") i.style.backgroundColor = "#1e1e1e";
    });
  } else {
    calc.classList.add("d-none");
    if (btnFlotante) btnFlotante.style.display = "block";
  }
};

window.pegarResultado = function () {
  const display = document.getElementById("displayCalc");
  if (display && ultimoInputEnfocado) {
    ultimoInputEnfocado.value = display.value;
    ultimoInputEnfocado.dispatchEvent(new Event("input"));
    window.addToCinta("-> Trasladado");
    ultimoInputEnfocado.focus();
  } else {
    window.mostrarAviso("Primero selecciona una celda");
  }
};

window.addToCinta = function (text) {
  const cinta = document.getElementById("cintaContable");
  if (!cinta) return;
  const line = document.createElement("div");
  line.style.color = text.includes("-") ? "#ff4444" : "inherit";
  line.style.fontSize = "0.75rem";
  line.innerText = text;
  cinta.appendChild(line);
  cinta.scrollTop = cinta.scrollHeight;
};

window.mostrarAviso = function (mensaje) {
  const toastEl = document.getElementById("toastAtajo");
  if (!toastEl) {
    console.log(mensaje);
    return;
  }
  const toastMensaje = document.getElementById("toastMensaje");
  if (toastMensaje) {
    toastMensaje.innerHTML = `<i class="fas fa-keyboard me-2"></i> ${mensaje}`;
  }
  new bootstrap.Toast(toastEl, { delay: 2000 }).show();
};

// ==========================================
// 3. LÓGICA DE NEGOCIO Y CÁLCULOS
// ==========================================

const formatearVisual = (input) => {
  let val = parseFloat(input.value.replace(/,/g, ""));
  if (!isNaN(val)) {
    input.value = val.toLocaleString("en-US", {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    });
  }
};

const calcularDiferencia = () => {
  let totalDebe = 0;
  let totalHaber = 0;
  let llenados = 0;
  const inputs = document.querySelectorAll(".input-contable");
  const diffValue = document.getElementById("diffValue");
  const cuadro = document.getElementById("cuadroCuadre");

  inputs.forEach((input) => {
    const nameInput = input.getAttribute("name") || "";
    const valorLimpio = input.value.replace(/,/g, "");
    const val = parseFloat(valorLimpio) || 0;

    if (nameInput === "220001" || nameInput === "220005") return;

    if (val !== 0) {
      llenados++;
      const saldo = input.getAttribute("data-saldo");
      if (["activo", "perdida", "costo_debe"].includes(saldo)) totalDebe += val;
      else if (["pasivo", "ganancia", "costo_haber"].includes(saldo))
        totalHaber += val;
    }
  });

  const diff = Math.abs(totalDebe - totalHaber);
  if (diffValue && cuadro) {
    diffValue.innerText = diff.toLocaleString("en-US", {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    });
    if (diff < 0.01 && llenados > 0) {
      diffValue.style.color = "#00ff00";
      diffValue.innerHTML = '<i class="fas fa-check-circle"></i> CUADRADO';
      if (!yaSonó) {
        sonidoExito.play().catch(() => {});
        yaSonó = true;
      }
    } else {
      diffValue.style.color = "#ff4444";
      yaSonó = false;
    }
  }
};

// ==========================================
// 4. EVENTOS Y DOM
// ==========================================

document.addEventListener("DOMContentLoaded", () => {
  const cuadro = document.getElementById("cuadroCuadre");
  if (cuadro) {
    document.body.appendChild(cuadro);
    cuadro.style.setProperty("right", "10px", "important");
    cuadro.style.setProperty("left", "auto", "important");
  }

  document.addEventListener("input", (e) => {
    if (e.target.classList.contains("input-contable")) {
      calcularDiferencia();
      const mappings = {
        111301: "220001",
        111303: "311001",
        111304: "311002",
        111305: "311001",
      };
      if (mappings[e.target.name]) {
        const target = document.querySelector(
          `[name="${mappings[e.target.name]}"]`,
        );
        if (target) {
          target.value = e.target.value;
          target.dispatchEvent(new Event("input"));
        }
      }
    }
  });

  document.addEventListener("focusin", (e) => {
    if (e.target.classList.contains("input-contable")) {
      ultimoInputEnfocado = e.target;
      e.target.value = e.target.value.replace(/,/g, "");
      e.target.select();
    }
  });

  document.addEventListener("focusout", (e) => {
    if (e.target.classList.contains("input-contable"))
      formatearVisual(e.target);
  });

  calcularDiferencia();
});

// ==========================================
// 5. ATAJOS DE TECLADO
// ==========================================

document.addEventListener(
  "keydown",
  (e) => {
    if (!e || !e.key) return; // Reparado: Movido arriba del todo para evitar fallos de lectura

    const isS = e.key === "s" || e.key === "S" || e.code === "KeiS";
    const key = e.key.toLowerCase();
    const isCtrl = e.ctrlKey || e.metaKey;
    const enInputTabla = e.target.classList.contains("input-contable");
    const calc = document.getElementById("calculadoraContable");
    const estaVisible = calc && !calc.classList.contains("d-none");

    if (isCtrl && isS) {
      e.preventDefault(); // DETIENE el "Guardar como" del navegador
      e.stopPropagation(); // EVITA que otros scripts lo vean
      console.log("💾 Iniciando guardado contable...");

      document.querySelectorAll(".input-contable").forEach((input) => {
        input.value = input.value.replace(/,/g, "");
      });

      const form = document.getElementById("balanceForm");
      const diffValue = document.getElementById("diffValue");
      const estaCuadrado =
        diffValue && diffValue.innerText.includes("CUADRADO");

      if (form) {
        if (estaCuadrado) {
          window.mostrarAviso("¡Cuadrado! Guardando datos...");
          form.submit();
        } else {
          if (
            confirm(
              "⚠️ La hoja NO está cuadrada. ¿Deseas guardar de todas formas?",
            )
          ) {
            form.submit();
          }
        }
      } else {
        console.error("No se encontró el formulario #balanceForm");
      }
      return false;
    }

    // Enter en tabla (Bloquea el envío automático en la última celda)
    if (key === "enter" && enInputTabla) {
      e.preventDefault();
      const inputs = Array.from(document.querySelectorAll(".input-contable"));
      const index = inputs.indexOf(e.target);

      if (index > -1 && index < inputs.length - 1) {
        inputs[index + 1].focus();
      } else {
        e.target.blur();
        console.log(
          "📝 Llegaste al final de la tabla. Usa Ctrl+S para guardar.",
        );
      }
      return;
    }

    // Alt + C
    if (e.altKey && key === "c") {
      e.preventDefault();
      window.toggleCalc();
      return;
    }

    // Teclas calculadora
    if (estaVisible && !enInputTabla) {
      if (/[0-9]/.test(key)) {
        e.preventDefault();
        window.addNum(key);
      }
      if (["+", "-", "*", "/"].includes(key)) {
        e.preventDefault();
        window["setOp"](key); // Usando corchetes el linter deja de buscar en el tipo estricto de Window
      }
      if (key === "." || key === ",") {
        e.preventDefault();
        window.addDecimal();
      }
      if (key === "backspace") {
        e.preventDefault();
        window.backspace();
      }
      if (key === "enter") {
        e.preventDefault();
        window.calcular();
      }
      if (key === "p") {
        e.preventDefault();
        window.pegarResultado();
      }
      if (key === "c" && !e.altKey) {
        e.preventDefault();
        window.clearCalc();
      }
      if (key === "escape") {
        e.preventDefault();
        window.toggleCalc();
      }
    }
  },
  { capture: true },
);

// Evitar que el formulario se envíe con un Enter accidental en cualquier lado
const formularioContable = document.getElementById("balanceForm");
if (formularioContable) {
  formularioContable.addEventListener("keydown", (e) => {
    if (e.key === "Enter") {
      if (e.target.tagName !== "BUTTON" && e.target.type !== "submit") {
        e.preventDefault();
      }
    }
  });
}

// ==========================================
// FUNCIÓN FISCAL
// ==========================================

window.calcularFiscal = function () {
  console.log("Calculando ISR y Reserva...");

  const baseElement = document.getElementById("baseContable");
  if (!baseElement) {
    console.error("No se encontró el elemento 'baseContable'");
    return;
  }

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
  if (document.getElementById("txtNeto"))
    document.getElementById("txtNeto").innerText = f.format(netoFinal);

  window.mostrarAviso("Cálculos fiscales actualizados");
};

// ==========================================
// FUNCIÓN DE AUDITORÍA (Sincronizada con Go)
// ==========================================

window.verAuditoria = function (nombre, detalleCuenta, interpretacion) {
  const modalElemento = document.getElementById("modalAuditoria");
  const modalTitulo = document.getElementById("modalAuditoriaTitulo");
  const modalCuerpo = document.getElementById("modalAuditoriaCuerpo");

  if (!modalElemento) {
    console.error(
      "No se encontró el modal con ID 'modalAuditoria' en el HTML.",
    );
    window.mostrarAviso("Error: No se encontró el modal de detalles.");
    return;
  }

  if (modalTitulo) {
    modalTitulo.innerHTML = `<i class="bi bi-shield-check me-2"></i> Auditoría: ${nombre}`;
  }

  if (modalCuerpo) {
    modalCuerpo.innerHTML = `
      <div class="p-2">
          <div class="mb-4">
              <h6 class="text-primary fw-bold mb-2">
                  <i class="bi bi-calculator me-2"></i> Fórmula y Valores (Origen Go):
              </h6>
              <div class="bg-dark text-success p-3 rounded font-monospace border shadow-sm" style="font-size: 0.95rem;">
                  ${detalleCuenta}
              </div>
          </div>
          <div class="mb-2">
              <h6 class="text-primary fw-bold mb-2">
                  <i class="bi bi-info-circle me-2"></i> Análisis Técnico:
              </h6>
              <p class="text-dark bg-light p-3 border-start border-4 border-primary rounded">
                  ${interpretacion}
              </p>
          </div>
          <div class="text-center mt-3">
              <small class="text-muted italic">
                  Los valores mostrados han sido extraídos directamente de los saldos de la Hoja de Trabajo.
              </small>
          </div>
      </div>
    `;
  }

  const miModal = new bootstrap.Modal(modalElemento);
  miModal.show();
};

// ==========================================
// 6. GENERACIÓN DE PDF REFORMADA
// ==========================================
const btnPdf = document.getElementById("botonpdf");
if (btnPdf) {
  btnPdf.addEventListener("click", function () {
    window.mostrarAviso("Generando Reporte Profesional por Páginas...");
    const fecha = new Date().toLocaleDateString("es-GT");
    const nombreArchivo = `Estados_Financieros_${fecha.replace(/\//g, "-")}.pdf`;
    const { jsPDF } = window.jspdf;
    const doc = new jsPDF("l", "pt", "letter");
    const totalPaginasReales = 5;

    const agregarPiePagina = (doc, numero) => {
      doc.setFontSize(9);
      doc.setTextColor(150);
      doc.text(
        `Página ${numero} de ${totalPaginasReales}`,
        doc.internal.pageSize.width - 80,
        doc.internal.pageSize.height - 30,
      );
    };

    const configurarTabla = (idSelector, titulo, yInicia) => {
      doc.setFontSize(14);
      doc.setTextColor(0, 0, 0);
      doc.text(titulo, 40, yInicia);
      doc.autoTable({
        html: idSelector,
        startY: yInicia + 20,
        theme: "grid",
        styles: {
          fontSize: 7.5,
          font: "courier",
          cellPadding: 2,
          fillColor: [255, 255, 255],
          textColor: [0, 0, 0],
        },
        headStyles: {
          fillColor: [230, 230, 230],
          textColor: [0, 0, 0],
          fontStyle: "bold",
        },
        didParseCell: function (data) {
          const el = data.cell.raw;
          if (el && el.tagName === "INPUT") data.cell.text = el.value;
          if (el && el.classList?.contains("border-top"))
            data.cell.styles.lineWidth = { top: 1.5 };
          if (el && el.classList?.contains("border-double")) {
            data.cell.styles.lineWidth = { top: 1, bottom: 2.5 };
            data.cell.styles.fontStyle = "bold";
          }
          if (data.column.index === 0) data.cell.styles.halign = "left";
          else if (data.column.index > 1) data.cell.styles.halign = "right";
        },
      });
    };

    // PÁGINA 1: HOJA DE TRABAJO
    configurarTabla("#hojaTrabajo", "HOJA DE TRABAJO - " + fecha, 50);
    agregarPiePagina(doc, 1);

    // PÁGINA 2: ESTADO DE RESULTADOS
    doc.addPage("letter", "p");
    configurarTabla("#resultados", "ESTADO DE RESULTADOS", 50);
    agregarPiePagina(doc, 2);

    // PÁGINA 3: BALANCE GENERAL
    doc.addPage("letter", "p");
    configurarTabla("#balance", "BALANCE GENERAL", 50);
    doc.autoTable({
      startY: doc.lastAutoTable.finalY + 40,
      theme: "plain",
      body: [
        ["f.__________________________", "f.__________________________"],
        ["Representante Legal", "Contador General"],
      ],
      styles: { halign: "center", fontSize: 9 },
    });
    agregarPiePagina(doc, 3);

    // PÁGINA 4: ANEXO CÁLCULO FISCAL
    doc.addPage("letter", "p");
    let finalY_BG = 50;
    doc.setFontSize(11);
    doc.setFont("helvetica", "bold");
    doc.text("ANEXO: CÁLCULO FISCAL", 40, finalY_BG);
    doc.autoTable({
      startY: finalY_BG + 10,
      margin: { left: 40 },
      tableWidth: 250,
      theme: "plain",
      body: [
        ["(+) ISR:", document.getElementById("txtISR")?.innerText || "0.00"],
        [
          "(+) Reserva Legal:",
          document.getElementById("txtReserva")?.innerText || "0.00",
        ],
        [
          "(=) Utilidad Neta:",
          document.getElementById("txtNeto")?.innerText || "0.00",
        ],
      ],
      styles: { fontSize: 9, font: "courier" },
    });
    agregarPiePagina(doc, 4); // Reparado: Ahora sí marca correlativamente la página 4

    // PÁGINA 5: ANÁLISIS FINANCIERO
    doc.addPage("letter", "p");
    doc.setFontSize(16);
    doc.setTextColor(0, 102, 204);
    doc.text("ANÁLISIS DE INDICADORES FINANCIEROS", 40, 50);

    const diccionarioFormulas = {
      "LIQUIDEZ CORRIENTE": "Activo Corriente / Pasivo Corriente",
      "ROTACIÓN DE INVENTARIO": "Costo de Ventas / Promedio Inventarios",
      "DÍAS DE INVENTARIO": "365 / Rotación de Inventarios",
      "NIVEL DE ENDEUDAMIENTO": "(Total Pasivo / Total Activo) * 100",
      "MARGEN DE UTILIDAD": "(Utilidad Neta / Ventas Netas) * 100",
      SOLVENCIA: "Total Activo / Total Pasivo",
    };

    const dataIndices = [];
    document.querySelectorAll(".card").forEach((card) => {
      const nombreRaw =
        card.querySelector(".card-subtitle")?.innerText.trim() || "";
      const valor = card.querySelector(".card-title")?.innerText.trim() || "";
      const interpretacion =
        card.querySelector(".card-text")?.innerText.trim() || "";
      const formula =
        diccionarioFormulas[nombreRaw.toUpperCase()] ||
        "Ver detalle en auditoría";

      if (nombreRaw && valor) {
        dataIndices.push([nombreRaw, valor, formula, interpretacion]);
      }
    });

    doc.autoTable({
      startY: 70,
      head: [["Indicador", "Valor", "Fórmula Aplicada", "Interpretación"]],
      body: dataIndices,
      theme: "striped",
      headStyles: { fillColor: [44, 62, 80], textColor: [255, 255, 255] },
      columnStyles: {
        0: { cellWidth: 100, fontStyle: "bold" },
        1: { cellWidth: 60, halign: "center", fontStyle: "bold" },
        2: { cellWidth: 120, fontSize: 8, font: "courier" },
        3: { cellWidth: "auto" },
      },
      styles: { fontSize: 8.5, cellPadding: 6 },
    });
    agregarPiePagina(doc, 5); // Reparado: Conteo correlativo final correcto

    doc.save(nombreArchivo);
    window.mostrarAviso("Reporte Maestro generado con éxito.");
  });
}
