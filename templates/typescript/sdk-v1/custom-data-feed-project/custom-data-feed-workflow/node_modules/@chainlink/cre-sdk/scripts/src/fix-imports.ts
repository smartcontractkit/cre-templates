import { glob } from 'fast-glob'
import { readFile, writeFile } from 'fs/promises'
import { dirname, join, relative } from 'path'

const fixImports = async () => {
	console.log('🔧 Fixing @cre/* imports in built files...')

	// Find all .js and .d.ts files in dist
	const files = await glob(['dist/**/*.js', 'dist/**/*.d.ts'], {
		cwd: process.cwd(),
		absolute: true,
	})

	let fixedCount = 0

	for (const file of files) {
		const content = await readFile(file, 'utf-8')

		// Replace @cre/* imports with relative paths
		const fixedContent = content.replace(/@cre\/([^'"]*)/g, (_, path) => {
			// Convert @cre/sdk/utils to relative path from current file
			const currentDir = dirname(file)
			const targetPath = join(process.cwd(), 'dist', path)
			const relativePath = relative(currentDir, targetPath)

			// Ensure it starts with ./ or ../
			return relativePath.startsWith('.') ? relativePath : `./${relativePath}`
		})

		if (fixedContent !== content) {
			await writeFile(file, fixedContent)
			fixedCount++
		}
	}

	console.log(`✅ Fixed imports in ${fixedCount} files`)
}

export const main = fixImports
