## **Caddy Defender Plugin**

The **Caddy Defender** plugin is a middleware for Caddy that allows you to block or manipulate requests based on the client's IP address. It is particularly useful for preventing unwanted traffic or polluting AI training data by returning garbage responses.

---

## **Features**

- **IP Range Filtering**: Block or manipulate requests from specific IP ranges.
- **Embedded IP Ranges**: Predefined IP ranges for popular AI services (e.g., OpenAI, DeepSeek, GitHub Copilot).
- **Custom IP Ranges**: Add your own IP ranges via Caddyfile configuration.
- **Multiple Responder Backends**:
  - **Block**: Return a `403 Forbidden` response.
  - **Garbage**: Return garbage data to pollute AI training.
  - **Custom**: Return a custom message.

---

## **Installation**

### **Using `xcaddy`**

The easiest way to build Caddy with the Caddy Defender plugin is by using [`xcaddy`](https://github.com/caddyserver/xcaddy), a tool for building custom Caddy binaries.

1. **Install `xcaddy`**:
   ```bash
   go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest
   ```

2. **Build Caddy with the Plugin**:
   Run the following command to build Caddy with the Caddy Defender plugin:
   ```bash
   xcaddy build --with github.com/jasonlovesdoggo/caddy-defender
   ```

   This will produce a `caddy` binary in the current directory.

3. **Run Caddy**:
   Use the built binary to run Caddy with your configuration:
   ```bash
   ./caddy run --config Caddyfile
   ```

---

### **Manual Build**

If you prefer to build the plugin manually:

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/jasonlovesdoggo/caddy-defender.git
   cd caddy-defender
   ```

2. **Build the Plugin**:
   ```bash
   go build -o caddy-defender
   ```

3. **Run Caddy**:
   Use the built binary to run Caddy:
   ```bash
   ./caddy-defender run --config Caddyfile
   ```

---

## **Configuration**

### **Caddyfile Syntax**

```caddyfile
defender <ip_ranges...> <responder> [responder_args...]
```

- `<ip_ranges...>`: A list of CIDR ranges to match against the client's IP.
- `<responder>`: The responder backend to use (`block`, `garbage`, or `custom`).
- `[responder_args...]`: Additional arguments for the responder backend (e.g., a custom message for the `custom` responder).

---

### **Examples**

#### **Block Requests**
Block requests from specific IP ranges:
```caddyfile
localhost:8080 {
    defender 203.0.113.0/24 198.51.100.0/24 block
    respond "Hello, world!"
}
```

#### **Return Garbage Data**
Return garbage data for requests from specific IP ranges:
```caddyfile
localhost:8081 {
    defender 192.168.1.0/24 garbage
    respond "Hello, world!"
}
```

#### **Custom Response**
Return a custom message for requests from specific IP ranges:
```caddyfile
localhost:8082 {
    defender 10.0.0.0/8 custom "Custom response message"
    respond "Hello, world!"
}
```

---

## **Embedded IP Ranges**

The plugin includes predefined IP ranges for popular AI services. These ranges are embedded in the binary and can be used without additional configuration.

| Service         | IP Ranges                                  |
|-----------------|--------------------------------------------|
| OpenAI          | [openai.go](ranges/fetchers/openai.go)     |
| DeepSeek        | [deepseek.go](ranges/fetchers/deepseek.go) |
| GitHub Copilot  | [github.go](ranges/fetchers/github.go)     |

More are welcome! for a precompiled list, see the [embedded results](ranges/data/generated.go)

---

## **GitHub Pages Documentation**

To host the documentation on GitHub Pages:

1. **Create a `docs` Folder**:
   Add a `docs` folder to your repository with an `index.md` file:
   ```
   caddy-defender/
   ├── docs/
   │   └── index.md
   └── ...
   ```

2. **Write the Documentation**:
   Copy the content of this README into `docs/index.md`.

3. **Enable GitHub Pages**:
   - Go to your repository's **Settings**.
   - Scroll down to the **Pages** section.
   - Set the source to the `docs` folder and save.

4. **Access the Documentation**:
   Your documentation will be available at:
   ```
   https://<username>.github.io/caddy-defender/
   ```

---

## **Contributing**

We welcome contributions! Here’s how you can get started:

### **Setting Up the Development Environment**

1. **Fork the Repository**:
   Fork the repository to your GitHub account.

2. **Clone the Repository**:
   ```bash
   git clone https://github.com/<your-username>/caddy-defender.git
   cd caddy-defender
   ```

3. **Build the Plugin**:
   Use `xcaddy` to build Caddy with the plugin:
   ```bash
   xcaddy build --with github.com/jasonlovesdoggo/caddy-defender
   ```

4. **Run Tests**:
   Add tests for your changes and run them:
   ```bash
   go test ./...
   ```

### **Submitting a Pull Request**

1. **Create a New Branch**:
   ```bash
   git checkout -b my-feature
   ```

2. **Commit Your Changes**:
   ```bash
   git commit -m "Add my feature"
   ```

3. **Push to GitHub**:
   ```bash
   git push origin my-feature
   ```

4. **Open a Pull Request**:
   Go to the repository on GitHub and open a pull request. Describe your changes and reference any related issues.

---

## **License**

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.

---

## **Acknowledgments**

- Inspired by the need to protect against unwanted AI traffic.
- Built with ❤️ using [Caddy](https://caddyserver.com).

---

Let me know if you need further assistance! 😊
