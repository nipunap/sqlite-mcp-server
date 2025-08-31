# Security Policy

## Supported Versions

The following versions of SQLite MCP Server are currently supported with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |
| < 0.1   | :x:                |

**Note**: As this is an early-stage project, we currently only support the latest minor version. Once we reach version 1.0, we will maintain security support for multiple versions.

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue in SQLite MCP Server, please report it responsibly.

### How to Report

**For security vulnerabilities, please do NOT create a public GitHub issue.**

Instead, please report security vulnerabilities through one of the following methods:

1. **GitHub Security Advisories** (Recommended)
   - Go to the [Security tab](https://github.com/nipunap/sqlite-mcp-server/security) of this repository
   - Click "Report a vulnerability"
   - Fill out the security advisory form

2. **Email** (Alternative)
   - Send an email with details to the repository maintainer
   - Include "SECURITY" in the subject line
   - Provide as much detail as possible about the vulnerability

### What to Include

When reporting a vulnerability, please include:

- **Description**: A clear description of the vulnerability
- **Steps to Reproduce**: Detailed steps to reproduce the issue
- **Impact**: What could an attacker accomplish with this vulnerability?
- **Affected Versions**: Which versions of the software are affected
- **Suggested Fix**: If you have ideas for how to fix the issue (optional)
- **Proof of Concept**: Code or screenshots demonstrating the vulnerability (if applicable)

### Response Timeline

We are committed to responding to security reports promptly:

- **Initial Response**: Within 48 hours of receiving the report
- **Status Updates**: Weekly updates on investigation progress
- **Resolution Timeline**: We aim to resolve critical vulnerabilities within 7 days, and other vulnerabilities within 30 days

### What to Expect

**If the vulnerability is accepted:**
- We will work with you to understand and reproduce the issue
- We will develop and test a fix
- We will coordinate the disclosure timeline with you
- We will credit you in the security advisory (unless you prefer to remain anonymous)
- We will release a security update and publish a security advisory

**If the vulnerability is declined:**
- We will provide a clear explanation of why we don't consider it a security issue
- We may suggest alternative ways to address your concerns
- You are welcome to discuss our decision if you disagree

## Security Best Practices

When using SQLite MCP Server, please follow these security best practices:

### Database Security
- **File Permissions**: Ensure database files have appropriate file system permissions
- **Access Control**: Limit access to database files to authorized users only
- **Backup Security**: Secure your database backups appropriately

### Network Security
- **Local Use**: SQLite MCP Server is designed for local use; avoid exposing it over networks
- **Input Validation**: Always validate and sanitize inputs when building applications on top of the server

### Configuration Security
- **Minimal Permissions**: Run the server with minimal required permissions
- **Regular Updates**: Keep the server updated to the latest supported version
- **Monitoring**: Monitor server logs for suspicious activity

## Scope

This security policy covers:
- The SQLite MCP Server core application
- Official Docker images (if any)
- Official documentation and examples

This policy does not cover:
- Third-party integrations or plugins
- User-created configurations or customizations
- Issues in dependencies (report these to the respective projects)

## Security Features

SQLite MCP Server includes the following security features:

- **Input Validation**: SQL injection protection through prepared statements
- **Error Handling**: Secure error messages that don't leak sensitive information
- **Resource Limits**: Protection against resource exhaustion attacks
- **Safe Defaults**: Secure default configuration settings

## Acknowledgments

We appreciate the security research community and will acknowledge researchers who report valid security vulnerabilities (unless they prefer to remain anonymous).

---

**Last Updated**: December 2024

For general questions about this security policy, please create a public GitHub issue with the "security" label.
