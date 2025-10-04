import crypto from 'node:crypto'
import fs from 'node:fs'
import os from 'node:os'
import path from 'node:path'
import zlib from 'node:zlib'

const exists = (p: string) => {
	try {
		fs.accessSync(p)
		return true
	} catch {
		return false
	}
}
const mkdirp = (p: string) => fs.mkdirSync(p, { recursive: true })

function computeAsset({ version }: { version: string }) {
	const platform = os.platform() // 'linux' | 'darwin' | 'win32'
	const arch = os.arch() // 'arm64' | 'x64' | ...

	const base = 'https://github.com/bytecodealliance/javy/releases/download'

	// macOS ARM64
	if (platform === 'darwin' && arch === 'arm64') {
		const name = `javy-arm-macos-${version}`
		return {
			gzUrl: `${base}/${version}/${name}.gz`,
			shaUrl: `${base}/${version}/${name}.gz.sha256`,
			cacheDir: path.join(os.homedir(), '.cache', 'javy', version, 'darwin-arm64'),
			outName: 'javy',
		}
	}

	// macOS x64
	if (platform === 'darwin' && arch === 'x64') {
		const name = `javy-x86_64-macos-${version}`
		return {
			gzUrl: `${base}/${version}/${name}.gz`,
			shaUrl: `${base}/${version}/${name}.gz.sha256`,
			cacheDir: path.join(os.homedir(), '.cache', 'javy', version, 'darwin-arm64'),
			outName: 'javy',
		}
	}

	// Linux ARM64
	if (platform === 'linux' && arch === 'arm64') {
		const name = `javy-arm-linux-${version}`
		return {
			gzUrl: `${base}/${version}/${name}.gz`,
			shaUrl: `${base}/${version}/${name}.gz.sha256`,
			cacheDir: path.join(os.homedir(), '.cache', 'javy', version, 'linux-arm64'),
			outName: 'javy',
		}
	}

	// Linux x64
	if (platform === 'linux' && arch === 'x64') {
		const name = `javy-x86_64-linux-${version}`
		return {
			gzUrl: `${base}/${version}/${name}.gz`,
			shaUrl: `${base}/${version}/${name}.gz.sha256`,
			cacheDir: path.join(os.homedir(), '.cache', 'javy', version, 'linux-arm64'),
			outName: 'javy',
		}
	}

	// Windows x64
	if (platform === 'win32' && arch === 'x64') {
		const name = `javy-x86_64-windows-${version}`
		return {
			gzUrl: `${base}/${version}/${name}.gz`,
			shaUrl: `${base}/${version}/${name}.gz.sha256`,
			cacheDir: path.join(os.homedir(), '.cache', 'javy', version, 'win32-x64'),
			outName: 'javy.exe',
		}
	}

	throw new Error(`Unsupported platform for Javy ${version}: ${platform}/${arch}`)
}

async function downloadText(url: string) {
	const res = await fetch(url)
	if (!res.ok) throw new Error(`HTTP ${res.status} for ${url}`)
	return res.text()
}

async function downloadBinary(url: string) {
	const res = await fetch(url)
	if (!res.ok) throw new Error(`HTTP ${res.status} for ${url}`)
	return Buffer.from(await res.arrayBuffer())
}

function verifySha256(buffer: Buffer, expectedHex: string) {
	const actual = crypto.createHash('sha256').update(buffer).digest('hex')
	if (!expectedHex || expectedHex.length < 64 || expectedHex !== actual) {
		throw new Error(`Checksum failed: expected ${expectedHex}, got ${actual}`)
	}
}

function writeFileAtomic(dest: string, data: Buffer, mode: number) {
	const tmp = `${dest}.${process.pid}.${Date.now()}.tmp`
	fs.writeFileSync(tmp, data, { mode })
	fs.renameSync(tmp, dest)
}

export async function ensureJavy({ version = 'v5.0.4' } = {}) {
	const { gzUrl, shaUrl, cacheDir, outName } = computeAsset({ version })
	const cacheBin = path.join(cacheDir, outName)
	if (exists(cacheBin)) return cacheBin

	mkdirp(cacheDir)

	const lock = path.join(cacheDir, '.lock')
	let weHoldLock = false
	try {
		const fd = fs.openSync(lock, fs.constants.O_CREAT | fs.constants.O_EXCL | fs.constants.O_WRONLY)
		fs.closeSync(fd)
		weHoldLock = true
	} catch {
		const start = Date.now()
		while (!exists(cacheBin) && Date.now() - start < 120000) {
			await new Promise((r) => setTimeout(r, 200))
		}
		if (exists(cacheBin)) return cacheBin
		try {
			fs.unlinkSync(lock)
		} catch {}
		const fd = fs.openSync(lock, fs.constants.O_CREAT | fs.constants.O_EXCL | fs.constants.O_WRONLY)
		fs.closeSync(fd)
		weHoldLock = true
	}

	try {
		const [shaText, gzBuf] = await Promise.all([downloadText(shaUrl), downloadBinary(gzUrl)])
		const expected = shaText.trim().split(/\s+/)[0]
		verifySha256(gzBuf, expected)
		const exeBuf = zlib.gunzipSync(gzBuf)
		writeFileAtomic(cacheBin, exeBuf, 0o755)
		return cacheBin
	} finally {
		if (weHoldLock) {
			try {
				fs.unlinkSync(lock)
			} catch {}
		}
	}
}
