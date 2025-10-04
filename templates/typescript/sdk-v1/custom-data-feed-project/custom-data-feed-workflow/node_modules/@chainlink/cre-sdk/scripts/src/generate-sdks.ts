import { file_capabilities_blockchain_evm_v1alpha_client } from '@cre/generated/capabilities/blockchain/evm/v1alpha/client_pb'
import { file_capabilities_internal_actionandtrigger_v1_action_and_trigger } from '@cre/generated/capabilities/internal/actionandtrigger/v1/action_and_trigger_pb'
import { file_capabilities_internal_basicaction_v1_basic_action } from '@cre/generated/capabilities/internal/basicaction/v1/basic_action_pb'
import { file_capabilities_internal_basictrigger_v1_basic_trigger } from '@cre/generated/capabilities/internal/basictrigger/v1/basic_trigger_pb'
import { file_capabilities_internal_consensus_v1alpha_consensus } from '@cre/generated/capabilities/internal/consensus/v1alpha/consensus_pb'
import { file_capabilities_internal_nodeaction_v1_node_action } from '@cre/generated/capabilities/internal/nodeaction/v1/node_action_pb'
import { file_capabilities_networking_http_v1alpha_client } from '@cre/generated/capabilities/networking/http/v1alpha/client_pb'
import { file_capabilities_networking_http_v1alpha_trigger } from '@cre/generated/capabilities/networking/http/v1alpha/trigger_pb'
import { file_capabilities_scheduler_cron_v1_trigger } from '@cre/generated/capabilities/scheduler/cron/v1/trigger_pb'
import { generateSdk } from '@cre/generator/generate-sdk'

export const main = () => {
	// Internal Testing SDKs
	generateSdk(file_capabilities_internal_basicaction_v1_basic_action, './src/generated-sdk')

	generateSdk(file_capabilities_internal_basictrigger_v1_basic_trigger, './src/generated-sdk')

	generateSdk(file_capabilities_internal_nodeaction_v1_node_action, './src/generated-sdk')

	generateSdk(file_capabilities_internal_consensus_v1alpha_consensus, './src/generated-sdk')

	generateSdk(
		file_capabilities_internal_actionandtrigger_v1_action_and_trigger,
		'./src/generated-sdk',
	)

	// Production SDKs

	generateSdk(file_capabilities_blockchain_evm_v1alpha_client, './src/generated-sdk')

	generateSdk(file_capabilities_networking_http_v1alpha_client, './src/generated-sdk')

	generateSdk(file_capabilities_networking_http_v1alpha_trigger, './src/generated-sdk')

	generateSdk(file_capabilities_scheduler_cron_v1_trigger, './src/generated-sdk')
}
