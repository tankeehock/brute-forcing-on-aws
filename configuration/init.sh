#!/bin/bash
aws s3 cp ${s3_user_data_uri} /home/ec2-user/app
chmod +x /home/ec2-user/app
/home/ec2-user/app -format "${target_access_key}" -n ${number_of_characters_for_brute_force} -secret "${target_secret_key}" -slave-index ${} -number-of-slaves ${number_of_slaves} -random -phone-number "${phone_number}"
