import fs from "fs";
import * as cdk from "aws-cdk-lib";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as iam from "aws-cdk-lib/aws-iam";
import { Construct } from "constructs";

interface RunnerConfig {
	readonly count: number;
	readonly instanceType: ec2.InstanceType;
	readonly amiId: string;
}

interface PisconStackProps extends cdk.StackProps {
	readonly runner: RunnerConfig;
	readonly sshPublicKeyPath: string;
}

export class AwsStack extends cdk.Stack {
	constructor(scope: Construct, id: string, props: PisconStackProps) {
		super(scope, id, props);

		const keyPair = new ec2.KeyPair(this, "ImportedKey", {
			publicKeyMaterial: fs.readFileSync(props.sshPublicKeyPath, "utf8"),
			keyPairName: "piscon-key-from-local",
		});

		const vpc = new ec2.Vpc(this, "PisconVPC", {
			ipAddresses: ec2.IpAddresses.cidr("10.0.0.0/24"),
			maxAzs: 1,
			natGateways: 0,
			subnetConfiguration: [
				{
					cidrMask: 24,
					name: "Public",
					subnetType: ec2.SubnetType.PUBLIC,
				},
			],
		});

		const portalSg = new ec2.SecurityGroup(this, "PortalSg", {
			vpc,
			description: "Security group for Portal",
		});

		const runnerSg = new ec2.SecurityGroup(this, "RunnerSg", {
			vpc,
			description: "Securty group for Runners",
		});

		portalSg.addIngressRule(
			runnerSg,
			ec2.Port.tcp(50051),
			"Allow gRPC from any Runner",
		);

		const portal = new ec2.Instance(this, "Portal", {
			vpc,
			instanceType: ec2.InstanceType.of(
				ec2.InstanceClass.T3A,
				ec2.InstanceSize.SMALL,
			),
			machineImage: ec2.MachineImage.genericLinux({
				"ap-northeast-1": "ami-0aec5ae807cea9ce0", // Ubuntu 24.04 x86_64
			}),
			securityGroup: portalSg,
			vpcSubnets: {
				subnetType: ec2.SubnetType.PUBLIC,
			},
			keyPair,
		});
		portal.connections.allowFrom(ec2.Peer.anyIpv4(), ec2.Port.SSH, "Allow SSH");
		portal.connections.allowFrom(
			ec2.Peer.anyIpv4(),
			ec2.Port.HTTPS,
			"Allow HTTPS",
		);
		portal.addToRolePolicy(
			new iam.PolicyStatement({
				effect: iam.Effect.ALLOW,
				actions: [
					"ec2:RunInstances",
					"ec2:TerminateInstances",
					"ec2:DescribeInstances",
					"ec2:DescribeTags",
					"ec2:CreateTags",
				],
				resources: ["*"],
			}),
		);

		new cdk.CfnOutput(this, "PortalPublicIp", {
			value: portal.instancePublicIp,
			description: "Portal server public IP address",
		});

		for (let i = 0; i < props.runner.count; i++) {
			const runner = new ec2.Instance(this, `Runner-${i}`, {
				vpc: vpc,
				instanceType: props.runner.instanceType,
				machineImage: ec2.MachineImage.genericLinux({
					"ap-northeast-1": props.runner.amiId,
				}),
				securityGroup: runnerSg,
				vpcSubnets: {
					subnetType: ec2.SubnetType.PUBLIC,
				},
				keyPair,
			});
			runner.connections.allowFrom(
				ec2.Peer.anyIpv4(),
				ec2.Port.SSH,
				"Allow SSH",
			);

			new cdk.CfnOutput(this, `RunnerIp-${i}`, {
				value: runner.instancePublicIp,
				description: `Runner ${i} public IP address`,
			});
		}
	}
}
