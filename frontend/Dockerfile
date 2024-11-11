# Usa una imagen base oficial de Node.js con la versión 18-alpine
FROM node:18-alpine

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /frontend

# Copiar los archivos package.json y package-lock.json (si existe)
COPY package*.json ./

# Instalar las dependencias
RUN npm install

# Copiar el resto de los archivos al contenedor
COPY . .

# Exponer el puerto en el que se ejecutará el frontend
EXPOSE 5173

# Comando para iniciar la aplicación
CMD ["npm", "run", "dev"]
