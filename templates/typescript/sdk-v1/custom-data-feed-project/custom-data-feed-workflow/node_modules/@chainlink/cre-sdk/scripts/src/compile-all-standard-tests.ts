// compile-all-standard-tests.ts
import path from 'node:path'
import fg from 'fast-glob'
import { main as compileWorkflow } from './compile-workflow'

/** Mirror the input file under .temp/<original path> and swap extension to .wasm */
export function toTempMirrorWasmPath(inputPath: string) {
	const resolved = path.resolve(inputPath)
	const relFromCwd = path.relative(process.cwd(), resolved)
	const pathParts = relFromCwd.split(path.sep)
	const filteredParts = pathParts.filter((part) => part !== 'src')
	const mirrored = path.join('.temp', ...filteredParts)
	return mirrored.replace(/\.[^.]+$/, '.wasm')
}

export const main = async (baseDir?: string, pattern?: string) => {
	// CLI fallback: bun test:standard:compile:all [baseDir] [pattern]
	const args = process.argv.slice(3)
	const root = baseDir ?? args[0] ?? 'src/standard_tests'
	const glob = pattern ?? args[1] ?? '**/test.ts'

	console.info('🚀 Compiling all standard tests...\n')
	console.info(`📂 Root:    ${path.resolve(root)}`)
	console.info(`🔎 Pattern: ${glob}\n`)

	// Find all test.ts files
	const files = await fg(`${root.replace(/\\/g, '/')}/${glob}`, { dot: false })

	if (files.length === 0) {
		console.error('❌ No standard test files found')
		process.exit(1)
	}

	console.info(`📋 Found ${files.length} test(s):`)
	for (const f of files) console.info(`   • ${f}`)
	console.info('')

	let success = 0
	let failed = 0
	const failures: string[] = []

	// Sequential for readable logs (switch to Promise.allSettled for parallel if desired)
	for (const testPath of files) {
		const rel = path.relative(process.cwd(), testPath)
		const wasmOut = toTempMirrorWasmPath(testPath)

		try {
			console.info(`🔨 Compiling: ${rel}`)
			console.info(`   → Output:  ${path.relative(process.cwd(), wasmOut)}`)
			await compileWorkflow(testPath, wasmOut)
			console.info(`✅ Success:   ${rel}\n`)
			success++
		} catch (err) {
			console.error(`❌ Failed:    ${rel}`)
			console.error(`   Error: ${err}\n`)
			failed++
			failures.push(rel)
		}
	}

	// Summary
	console.info('📊 Compilation Summary:')
	console.info(`   ✅ Successful: ${success}`)
	console.info(`   ❌ Failed:     ${failed}`)

	if (failures.length > 0) {
		console.info('   Failed tests:')
		for (const f of failures) console.info(`   • ${f}`)
		process.exit(1)
	}

	console.info(`\n🎉 All ${success} standard tests compiled successfully!`)
}

// Allow direct CLI usage
if (import.meta.main) {
	main().catch((e) => {
		console.error(e)
		process.exit(1)
	})
}
