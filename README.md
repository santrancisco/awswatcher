### awswatcher

While studying for my AWS exam, I deciced to put cloudformation + Alarm metric and some go into practice. 

awswatcher is a simple golang lambda function deployed by Cloudformation template through AWS Console.

The awswatcher lambda function responsibles for sending a slack notification to a specific channel via the "SLACK_WEBHOOK" environment variable attached to the function (can be set during cloudformation stack installation)

The `CloudWatch_Alarms_for_CloudTrail_API_Activity.json` cloudformation template is a modified version of the template you can find [here](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/use-cloudformation-template-to-create-cloudwatch-alarms.html). 

The original template includes an SNSTopic with ability to send email notification to an user. I modified the template and included the lambda function, the role required for the function.

Cloudformation template is, in my opinion, best fit for this job over terraform and there is a reason for that which you will find in a blog post i have [here](https://blog.ebfe.pw)

### Prerequisites

In order for cloudformation to run and pull the lambda function binary, you will need to first create an S3 bucket in the same region as where you currently have the Cloudwatch Log being pushed to from CloudTrail. This is also where you want to run your Cloudformation template. So to make it simple, here is the list of thigns you need to do before executing this cloudfomration stack:

 - Create or reuse the Cloudtrail to Cloudwatch Log integration and note down the cloudwatch log group name.

 - Create an S3 bucket in the same region with Cloudwatch you setup above. This bucket will be used to deploy the code to our lambda function
 
 - Run `make` to create awswatcher.zip under ./bin folder.

 - Uplpad the awswatcher.zip file to the S3 bucket.

 - Upload this cloudformation template and fill in the information with what we have above.

 - Cloudformation will create our awesome stack with all the metrics, alarm, lambda function + SNS Topic.

 - Have a coffee and wait for everything to work magically
