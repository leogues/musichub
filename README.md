<h1 align="center" id="readme-top">
  <img alt="musichub" title="MusicHub" src="https://github.com/leogues/musichub/blob/dev/frontend/src/assets/images/title.png?raw=true" width="500px">
</h1>

<p align="center">
  <img alt="GitHub language count" src="https://img.shields.io/github/languages/count/leogues/musichub?color=%2304D361">

  <img alt="License" src="https://img.shields.io/badge/license-MIT-%2304D361">

  <a href="https://github.com/leogues/musichub/stargazers">
    <img alt="Stargazers" src="https://img.shields.io/github/stars/leogues/musichub?style=social">
  </a>
</p>

## Sobre o projeto

[![demo][demo-image]](https://music.leogues.com.br/)

[demo-image]: https://github.com/leogues/musichub/blob/dev/frontend/src/assets/images/demo.png?raw=true

O Music Hub é uma plataforma que sincroniza sua biblioteca de músicas com diversos serviços de streaming. Ele permite criar smart links e transferir suas músicas entre diferentes serviços. Ainda está em construção.

## 🚀 Começando

Este guia fornecerá instruções detalhadas sobre como configurar e rodar o projeto "MusicHub" localmente. Siga os passos abaixo para ter uma cópia local instalada e funcionando.

### Pre-requisites

Antes de começar, certifique-se de ter os seguintes softwares instalados:

1. NPM: O gerenciador de pacotes do Node.js.

```sh
npm install npm@latest -g
```

2. Docker: Plataforma para desenvolvimento e execução de aplicações em containers.

[Docker](https://www.docker.com/)

3. Angular CLI:

```sh
npm install -g @angular/cli@17.3
```

4. Configure as variaveis de desenvolvimento:

Crie o arquivo .env com base no .env.example do backend

### Instalação

Siga os passos abaixo para configurar o ambiente de desenvolvimento local:

1. Clone o repositorio

```sh
git clone https://github.com/leogues/live-link.git
```

2. Instale NPM packages

```sh
npm install
```

3. Execute o docker compose

```sh
 docker compose -f docker-compose.dev.yml up -d
```

<p align="right">(<a href="#readme-top">Voltar para o topo</a>)</p>
