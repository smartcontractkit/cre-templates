#!/usr/bin/env bun

import { spawn } from 'node:child_process'
import { copyFileSync, mkdirSync } from 'node:fs'
import { join } from 'node:path'

const builtWasmPath = join(
	process.cwd(),
	'src',
	'javy_chainlink_sdk',
	'target',
	'wasm32-wasip1',
	'release',
	'javy_chainlink_sdk.wasm',
)
const distWasmPath = join(process.cwd(), 'dist', 'javy_chainlink_sdk.wasm')
const witFilePath = join(process.cwd(), 'src', 'workflow.wit')
const distWitFilePath = join(process.cwd(), 'dist', 'workflow.wit')

export const main = async () => {
	const pluginDir = join(process.cwd(), 'src', 'javy_chainlink_sdk')

	console.info('\n\n---> Compiling Chainlink SDK Javy plugin (Rust) \n\n')

	return new Promise<void>((resolve, reject) => {
		const buildProcess = spawn('cargo', ['build', '--target', 'wasm32-wasip1', '--release'], {
			cwd: pluginDir,
			stdio: 'inherit',
			shell: true,
		})

		buildProcess.on('close', (code) => {
			if (code === 0) {
				mkdirSync('dist', { recursive: true })
				copyFileSync(builtWasmPath, distWasmPath)
				copyFileSync(witFilePath, distWitFilePath)

				console.info('✅ Done!')
				resolve()
			} else {
				console.error(`❌ Plugin build failed with code ${code}`)
				reject(new Error(`Plugin build failed with code ${code}`))
			}
		})

		buildProcess.on('error', (error) => {
			console.error('❌ Failed to start build process:', error)
			reject(error)
		})
	})
}

main()
