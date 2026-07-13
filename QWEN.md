# States

[![CI](https://github.com/<OWNER>/<REPO>/actions/workflows/ci.yml/badge.svg)](https://github.com/<OWNER>/<REPO>/actions/workflows/ci.yml)

## Descripción

States es una aplicación web construida en Go con Fiber y Handlebars, diseñada para gestionar usuarios, suscripciones, soporte y generación de estados financieros / auditorías de costos.

Incluye:
- Autenticación JWT basada en cookies.
- Registro de usuarios con verificación por correo.
- Cálculo de estados financieros y costos industriales.
- Panel de usuario con historial de transacciones.
- Panel administrativo con gestión de soporte y usuarios.
- Integración con PostgreSQL y migraciones automáticas.
- Despliegue con Docker y Nginx.

## Características principales

- Registro de usuarios con validación de captchas.
- Envío de correos de verificación usando Resend.
- Login seguro con JWT en cookie `jwt`.
- Rutas públicas, protegidas y administrativas.
- Generación de estados financieros (EEFF) para empresas comerciales.
- Generación de cálculos de costos industriales.
- Soporte técnico con mensajería interna.
- Webhook para procesar notificaciones de pago de pasarela QPayPro.
- Carga dinámica de nomenclatura contable desde la tabla `nomenclatura`.

## Estructura del proyecto

- `main.go` - punto de entrada, config de Fiber, carga de DB, rutas y middleware.
- `db/` - conexión a PostgreSQL y migración de esquemas.
- `config/` - carga y clasificación de cuentas/nomenclatura contable.
- `middleware/` - autenticación, autorización y estado de sesión.
- `ctrl/` - controladores HTTP para login, registro, perfiles, soporte y reports.
- `models/` - entidades GORM y estructuras para los cálculos financieros.
- `servicios/` - integración con servicios externos (correo de verificación).
- `utils/` - helpers de plantilla Handlebars.
- `views/` - plantillas Handlebars (.hbs).
- `estaticos/` - archivos estáticos.
- `nginx/` - configuración de proxy reverso para producción.
- `Dockerfile`, `docker-compose.yml`, `docker-compose.prod.yml` - contenedores de desarrollo y producción.

## Requisitos

- Go 1.25.4
- PostgreSQL como base de datos
- Docker y Docker Compose (opcional, pero recomendado para desarrollo y despliegue)

## Variables de entorno

Crea un archivo `.env` con las siguientes variables:

```env
PORT=3000
APP_ENV=production
DB_HOST=...
DB_USER=...
DB_PASSWORD=...
DB_NAME=...
DB_PORT=...
DB_SSL=...
JWT_SECRET=...
EEFFS_APP=...
RESEND_API_KEY=...
```

### Descripción de variables

- `PORT`: puerto en el que se levanta la app.
- `APP_ENV`: entorno de ejecución (`development`, `production`, etc.).
- `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`, `DB_SSL`: conexión a PostgreSQL.
- `JWT_SECRET`: clave para firmar y validar tokens JWT.
- `EEFFS_APP`: token secreto de reCAPTCHA usado en el registro.
- `RESEND_API_KEY`: API Key para enviar correos de verificación con Resend.

> Nota: no subas este archivo con credenciales reales a un repositorio público.

## Instalación local

1. Clona el repositorio y entra en la carpeta:

```bash
git clone <repositorio>
cd States
```

2. Copia el ejemplo de variables de entorno y ajusta los valores:

```bash
cp .env .env.local
```

3. Descarga dependencias de Go:

```bash
go mod download
```

4. Ejecuta la aplicación:

```bash
go run main.go
```

5. Abre `http://localhost:3000` en tu navegador.

## Desarrollo con Docker

### Modo desarrollo

Este proyecto incluye un `Dockerfile` de desarrollo que instala `air` para recarga automática.

```bash
docker compose up --build
```

### Modo producción

Para producción se usa `Dockerfile.prod` y `docker-compose.prod.yml` con Nginx:

```bash
docker compose -f docker-compose.prod.yml up --build
```

## Rutas importantes

### Públicas

- `GET /` - página de inicio.
- `GET /loguin` - formulario de inicio de sesión.
- `POST /loguin` - envío de credenciales.
- `GET /register` - formulario de registro.
- `POST /register` - registro de usuario.
- `GET /verificar` - verificación de cuenta por token.
- `GET /about` - página de información.
- `GET /manual` - manual de usuario.

### Protegidas (requieren sesión)

- `GET /eeff` - módulo de hoja de trabajo comercial.
- `POST /estados` - generación de estados financieros.
- `GET /perfil` - perfil de usuario y transacciones.
- `GET /costosform` - módulo de costos industriales.
- `POST /costos` - cálculo de costos de producción.
- `GET /soport` - soporte técnico.
- `POST /soporte/enviar` - envío de consulta técnica.
- `GET /logout` - cierre de sesión.

### API / admin

- `POST /api/pagos/qpay-webhook` - webhook de pagos QPayPro.
- `GET /api/admin/soporte` - lista de consultas pendientes (admin).
- `POST /api/admin/soporte/responder/:id` - responder una consulta.
- `GET /api/admin/dashboard` - dashboard administrativo.
- `GET /api/admin/usuario/:id` - obtener detalles de un usuario.
- `GET /api/admin/crearcuenta` - panel de creación de cuenta.
- `POST /api/admin/crearcuenta` - crear una cuenta contable.
- `GET /api/admin/eliminar-cuenta/:codigo` - eliminar cuenta contable.
- `POST /api/admin/editar-cuenta` - editar cuenta contable.

## Base de datos y migraciones

La aplicación usa GORM para migrar automáticamente estas entidades al iniciar:

- `models.User`
- `models.Mensaje`
- `models.CueDB` (tabla `nomenclatura`)
- `models.Transaccion`

También carga la nomenclatura contable desde la tabla `nomenclatura` para usarla en el módulo financiero.

## Seguridad y recomendaciones

- La cookie de sesión se llama `jwt` y está configurada como `HTTPOnly` y `Secure`.
- Asegúrate de usar HTTPS en producción para que `Secure=true` funcione correctamente.
- Revisa que `JWT_SECRET`, `RESEND_API_KEY` y `EEFFS_APP` estén correctamente configurados.
- Cambia `loguin` si prefieres normalizar la ruta a `login` en vistas y rutas.

## Observaciones adicionales

- El proyecto está preparado para traducir la contabilidad en estados financieros y costos industriales.
- La carga de nomenclatura se realiza desde la DB para mantener actualizada la lista de cuentas.
- El controlador de correo usa Resend para envío de verificación.

## Próximos pasos sugeridos

- Añadir pruebas unitarias e integración.
- Agregar un `README` de variables de configuración y desplegar en un entorno controlado.
- Validar los endpoints de admin y el webhook de pagos en un entorno de pruebas.
- Mejorar el error handling en los formularios y el flujo de suscripción.

## CI / GitHub Actions

Se añadió un workflow de CI que ejecuta `go test ./... -v` en `push` y `pull_request`.

- Archivo: `.github/workflows/ci.yml`
- Acción: ejecuta tests de Go en `ubuntu-latest` usando Go `1.25.4`.

Actualiza el badge arriba reemplazando `<OWNER>/<REPO>` por tu organización y repositorio reales.

## GitHub Secrets recomendados

Para que el CI y el despliegue funcionen correctamente (y para mantener seguros tus datos), añade estos secretos en `Settings > Secrets and variables > Actions` del repositorio:

- `DB_HOST` — Host de la base de datos (ej. `db.example.com`)
- `DB_USER` — Usuario de la base de datos
- `DB_PASSWORD` — Contraseña de la base de datos
- `DB_NAME` — Nombre de la base de datos
- `DB_PORT` — Puerto de la base de datos
- `DB_SSL` — `require` / `disable` según tu entorno
- `JWT_SECRET` — Clave para firmar JWT
- `RESEND_API_KEY` — API key para Resend (envío de correos)
- `EEFFS_APP` — Clave de reCAPTCHA (o `RECAPTCHA_SECRET` si prefieres renombrar)

## 🎨 Criterios de Diseño, Código y Lenguaje
- **Minimalismo:** Buscamos la menor cantidad de líneas de código posibles. La simplicidad es la base; código limpio y directo.
- **Frontend:** Estilos basados estrictamente en Bootstrap (estructurado por etiquetas, grid y flexbox). Evitar CSS inline innecesario y mantener las plantillas `.hbs` limpias.
- **Lenguaje** Utiliza español para los comentarios variables y para escribirme.
