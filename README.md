# Go AWS Profile Switch

Simple cli utility, switch your default profile easily!

# Install
```
go get -u github.com/okamos/go-aws-profile-switch/cmd/awsswitch
```

# Usage
```
Usage: awsswitch [--version] [--help] <command> [<args>]

Available commands are:
    ls    Lists available AWS profiles
    sw    Switches default your AWS profile
```

# Examples
Given the following `~/.aws/credentials` file.

```
[default]
aws_access_key_id=profile1_id
aws_secret_access_key=profile1_secret
region=us-east-1
output=json

[profile1]
aws_access_key_id=profile1_id
aws_secret_access_key=profile1_secret
region=us-east-1
output=json

[profile2]
aws_access_key_id=profile2_id
aws_secret_access_key=profile2_secret
region=us-west-2
output=json
```

`awsswitch ls` output is

```
Available profiles
  default
* profile1
  profile2
```

`awsswitch sw -p profile2` will be Switching default profile to profile2,  
and output is

```
Your default profile is overwrote
> aws_access_key_id=profile2_id
> aws_secret_access_key=profile2_secret
> region=us-west-2
> output=json

From
> aws_access_key_id=profile1_id
> aws_secret_access_key=profile1_secret
> region=us-west-1
> output=json
```

Now `awsswitch ls` output is

```
Available profiles
  default
  profile1
* profile2
```

Then `awsswitch sw -p profile2` isn't rewrite default profile, because default  
profile is profile2 already.

```
Your default profile set profile2 already
```

# How it works
Global rules  

1. Triming all leading and trailing spaces.

## ls command rules
1. Obtains all profiles.
2. Add `*` to profile name, if exists default profile and any profile has a same aws_access_key_id.
3. Display profiles.

## sw command rules
1. Obtains all profiles.
2. When start with `#` character, Regarded comment line.
3. Your input profile is set already, awsswitch isn't rewrite profile.
4. When not found the default aws_access_key_id from other profiles, awsswitch isn't rewrite default profile.
5. Write default profile and other profiles, comments to `~/.aws/credentials`.
