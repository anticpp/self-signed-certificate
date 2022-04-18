/*
Application behavior is defined by it's program(binary/script) and it's configurations. Thus, how to handle configurations is an important job for most applications. In fact there are different config sources, for example, config files, commandline arguments, or environments.
It's not so easy jobs to handle these different config sources efficiently, you have to write different codes for different config source.
More over, we always specify some part of configurations with config file, the other part with commandline arguments. But there is no way to specify a single configuration with multiple config sources, or say, it's another complicated job to handle it.

Package config try to provide an unify way for different config sources, you can specify your configurations with any type of config source as you like, just consider your deployment senario. For example, if you just want to make a local test run, you can use commandline, the easiest way. If you want to deploy in Docker container, sometimes, use environments will be easier. If you need to make a production deployment, use config files, will be better.
*/
package config
