#!/usr/bin/env node
import "dotenv/config";
import * as cdk from "aws-cdk-lib";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import { AwsStack } from "../lib/aws-stack";
import path from "path";

const app = new cdk.App();

const runnerAmiId = process.env.RUNNER_AMI_ID;
if (!runnerAmiId) {
	throw new Error("RUNNER_AMI_ID is not configured");
}

const config = {
	runner: {
		count: parseInt(process.env.RUNNER_COUNT || "1"),
		instanceType: new ec2.InstanceType(process.env.RUNNER_TYPE || "t3a.small"),
		amiId: runnerAmiId,
	},
	sshPublicKeyPath:
		process.env.SSH_PUBLIC_KEY ||
		path.join(process.env.HOME!, ".ssh/id_ed25519.pub"),
};

new AwsStack(app, "PisconStack", {
	runner: config.runner,
	sshPublicKeyPath: config.sshPublicKeyPath,
	env: {
		account: process.env.CDK_DEFAULT_ACCOUNT,
		region: process.env.CDK_DEFAULT_REGION || "ap-northeast-1",
	},
});
