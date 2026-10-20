import {
  CronCapability,
  NITRO_REGIONS,
  Runner,
  handlerInTee,
  type TeeConstraint,
  type TeeRuntime,
  type Workflow,
} from "@chainlink/cre-sdk";

export type Config = {
  schedule: string;
};

// Run in an AWS Nitro enclave in us-west-2 (staging confidential-workflows
// region). Update this when other TEE types or regions become available.
const teeRequirements: TeeConstraint = [
  { tee: "nitro", regions: [...NITRO_REGIONS] },
];

// onCronTrigger runs inside the enclave. The SDK routes Math.random() to the
// TEE-local seed for handlers registered with handlerInTee, so each execution
// gets a fresh value that differs across runs and concurrent executions.
export const onCronTrigger = (runtime: TeeRuntime<Config>): string => {
  const value = Math.floor(Math.random() * 0x1_0000_0000) >>> 0;
  runtime.log(`random-number-generated random_value=${value}`);
  return value.toString();
};

export const initWorkflow = (config: Config): Workflow<Config> => {
  if (!config.schedule) {
    throw new Error("config requires schedule");
  }

  const cron = new CronCapability();

  return [
    handlerInTee(
      cron.trigger({ schedule: config.schedule }),
      onCronTrigger,
      teeRequirements,
    ),
  ];
};

export async function main() {
  const runner = await Runner.newRunner<Config>();
  await runner.run(initWorkflow);
}
