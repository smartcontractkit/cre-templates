export const MintingConsumer = [
	{
		type: 'constructor',
		inputs: [
			{ name: '_stablecoin', type: 'address', internalType: 'address' },
			{ name: '_expectedAuthor', type: 'address', internalType: 'address' },
			{ name: '_expectedWorkflowName', type: 'bytes10', internalType: 'bytes10' },
		],
		stateMutability: 'nonpayable',
	},
	{
		type: 'function',
		name: 'EXPECTED_AUTHOR',
		inputs: [],
		outputs: [{ name: '', type: 'address', internalType: 'address' }],
		stateMutability: 'view',
	},
	{
		type: 'function',
		name: 'EXPECTED_WORKFLOW_NAME',
		inputs: [],
		outputs: [{ name: '', type: 'bytes10', internalType: 'bytes10' }],
		stateMutability: 'view',
	},
	{
		type: 'function',
		name: 'onReport',
		inputs: [
			{ name: 'metadata', type: 'bytes', internalType: 'bytes' },
			{ name: 'report', type: 'bytes', internalType: 'bytes' },
		],
		outputs: [],
		stateMutability: 'nonpayable',
	},
	{
		type: 'function',
		name: 'stablecoin',
		inputs: [],
		outputs: [{ name: '', type: 'address', internalType: 'contract StablecoinERC20' }],
		stateMutability: 'view',
	},
	{
		type: 'function',
		name: 'supportsInterface',
		inputs: [{ name: 'interfaceId', type: 'bytes4', internalType: 'bytes4' }],
		outputs: [{ name: '', type: 'bool', internalType: 'bool' }],
		stateMutability: 'pure',
	},
	{
		type: 'event',
		name: 'MintExecuted',
		inputs: [
			{ name: 'recipient', type: 'address', indexed: true, internalType: 'address' },
			{ name: 'amount', type: 'uint256', indexed: false, internalType: 'uint256' },
			{ name: 'bankReference', type: 'bytes32', indexed: true, internalType: 'bytes32' },
		],
		anonymous: false,
	},
	{
		type: 'event',
		name: 'MintInstructionReceived',
		inputs: [
			{ name: 'recipient', type: 'address', indexed: true, internalType: 'address' },
			{ name: 'amount', type: 'uint256', indexed: false, internalType: 'uint256' },
			{ name: 'bankReference', type: 'bytes32', indexed: true, internalType: 'bytes32' },
			{ name: 'timestamp', type: 'uint256', indexed: false, internalType: 'uint256' },
		],
		anonymous: false,
	},
	{
		type: 'error',
		name: 'InvalidAuthor',
		inputs: [
			{ name: 'received', type: 'address', internalType: 'address' },
			{ name: 'expected', type: 'address', internalType: 'address' },
		],
	},
	{
		type: 'error',
		name: 'InvalidWorkflowName',
		inputs: [
			{ name: 'received', type: 'bytes10', internalType: 'bytes10' },
			{ name: 'expected', type: 'bytes10', internalType: 'bytes10' },
		],
	},
	{ type: 'error', name: 'MintFailed', inputs: [] },
] as const

