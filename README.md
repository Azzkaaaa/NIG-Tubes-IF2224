# NIG - TUBES - IF2224

<div align="center">
  <img width="100%" src="https://capsule-render.vercel.app/api?type=waving&height=300&color=timeGradient&text=NIG%20K-03%20G-02&reversal=true&fontColor=ffffff&animation=twinkling&stroke=fffffff&strokeWidth=4&fontSize=0" />
</div>

<p align="center">
  <img src="https://img.shields.io/badge/Status-Done-008000" />
  <img src="https://img.shields.io/badge/Version-1.0.0-brightgreen" />
  <img src="https://img.shields.io/badge/License-MIT-yellowgreen" />
  <img src="https://img.shields.io/badge/Built_With-Go-blue" />
</p>

<h1 align="center">
  <img src="https://readme-typing-svg.herokuapp.com?font=Fira+Code&pause=500&color=81a1c1&center=true&vCenter=true&width=600&lines=13523123,+13523137,+13523161,+13523162,+13523163;Bimo,+Azka,+Arlow,+Riza,+dan+Filbert" alt="R.Bimo, Rafael, dan Frederiko" />
</h1>

## Overview

Sebuah Interpreter bahasa Pascal.

## Requirements

- Punya Go-lang

## Cara Install Go-lang

### Windows

1.  **Download the Installer:**
    * Go to the official Go downloads page: https://go.dev/dl/
    * Download the latest stable version's **MSI installer** for Windows (x64).

2.  **Run the Installer:**
    * Double-click the downloaded **.msi** file.
    * Follow the installation wizard prompts. The installer automatically places Go in a standard location (like C:\Program Files\Go or C:\Go) and updates your **PATH** environment variable.

3.  **Verify the Installation:**
    * Open a **new** Command Prompt or PowerShell window.
    * Run the command:
        go version
    * You should see the installed Go version (e.g., go version go1.x.x windows/amd64).

---

### Linux (with Package Manager)

This method ensures you get the latest version of Go.

* **For APT (Debian/Ubuntu):**
    ```terminal
    sudo apt update
    sudo apt install golang-go
    ```

* **For dnf (Fedora/CentOS/RHEL):**
    ```terminal
    sudo dnf install golang
    ```

* **For Pacman (Arch based):**
    ```terminal
    sudo pacman -Syu
    sudo pacman -S go
    go version
    ```

* **For apk (Alpine):**
    ```terminal
    apk update
    apk add go
    go version
    ```

## Cara run Lexer
1. pada root repository ketik:

    windows:
    ```bash
    go build -o bin/pslex.exe ./cmd/pslex
    ``` 

    linux:
    ```bash
    go build -o bin/pslex ./cmd/pslex
    ``` 

2. Lalu ketik:

   windows:
   ```bash
   ./bin/pspar.exe --rules config/tokenizer_indo.json --input test/milestone-2/(nama file pascal.pas)
   ```

   linux:
   ```bash
   ./bin/pspar --rules config/tokenizer_indo.json --input test/milestone-2/(nama file pascal.pas)
   ```

contoh:
./bin/pslex.exe --rules src/rules/tokenizer.json --input test/milestone-2/test1.pas

## Cara run Parser
1. pada root repository ketik:
    
    windows:
    ```bash
    go build -C psc/cmd/pspar -o ../../../bin/pspar.exe .
    ```

    linux:
   ```bash
   go build -C psc/cmd/pspar -o ../../../bin/pspar .
   ```

3. Lalu ketik:

    windows:
    ```bash
    ./bin/pspar.exe --rules config/tokenizer_indo.json --input test/milestone-2/(nama file pascal.pas)
    ```

   linux:
   ```bash
   ./bin/pspar --rules config/tokenizer_indo.json --input test/milestone-2/(nama file pascal.pas)
   ```

## Pembagian Tugas

<div>
    <table align="center">
    <tr>
        <th align="center">User</th>
        <th align="center">Job</th>
    </tr>
    <tr>
        <td align="center">
        <a href="https://github.com/Cola1000">
            <img src="https://avatars.githubusercontent.com/u/143616767?v=4" width="80px" style="border-radius: 50%;" alt="Cola1000"/><br />
            <sub><b>Rhio Bimo Prakoso S</b></sub>
        </a>
        </td>
        <td align="center">
          M1: Fotografer, Konsumsi, Femboy Generator, Testing, dan Laporan
        <br>
          M2: Tester, Laporan, Code
        </td>
    </tr>
    <tr>
        <td align="center">
        <a href="https://github.com/Azzkaaaa">
            <img src="https://avatars.githubusercontent.com/u/167986692?v=4" width="80px" style="border-radius: 50%;" alt="V-Kleio"/><br />
            <sub><b>M Aulia Azka</b></sub>
        </a>
        </td>
        <td align="center">
          M1: Engine dan Testing
        <br>
          M2: Laporan, Tester
        </td>
    </tr>
    <tr>
        <td align="center">
        <a href="https://github.com/Arlow5761">
            <img src="https://avatars.githubusercontent.com/u/96019562?v=4" width="80px" style="border-radius: 50%;" alt="susTuna"/><br />
            <sub><b>Arlow Emmanuel Hergara</b></sub>
        </a>
        </td>
        <td align="center">
          M1: Lexer (Graph and Code)
        <br>
          M2: Laporan, Code
        </td>
    </tr>
        <tr>
        <td align="center">
        <a href="https://github.com/L4mbads">
            <img src="https://avatars.githubusercontent.com/u/85736842?v=4" width="80px" style="border-radius: 50%;" alt="susTuna"/><br />
            <sub><b>Fahcriza Ahmad Setiyono</b></sub>
        </a>
        </td>
        <td align="center">
          M1: Engine dan Testing
        <br>
          M2: Laporan, Code
        </td>
    </tr>
        <tr>
        <td align="center">
        <a href="https://github.com/filbertengyo">
            <img src="https://avatars.githubusercontent.com/u/163801345?v=4" width="80px" style="border-radius: 50%;" alt="susTuna"/><br />
            <sub><b>Filbert Engyo</b></sub>
        </a>
        </td>
        <td align="center">
          M1: Engine dan Testing
        <br>
          M2: Tester, Laporan, Code
        </td>
    </tr>
    </table>
</div>

<div align="center" style="color:#6A994E;"> 🌿 Please Donate for Charity! 🌿</div>

<p align="center">
  <a href="https://tiltify.com/@cdawg-va/cdawgva-cyclethon-4" target="_blank">
    <img src="https://assets.tiltify.com/uploads/cause/avatar/4569/blob-9169ab7d-a78f-4373-8601-d1999ede3a8d.png" alt="IDF" style="height: 80px;padding: 20px" />
  </a>
</p>
