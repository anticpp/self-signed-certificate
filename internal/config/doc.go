/*
It's very common situaion for application handling configurations, these configurations may be  read from config file, commandline, or environments.
How to handle these different config sources is not so easy, becasuse you have to write different codes for each different config source.
And more, we always split our configuratrions to different kinds, some kind be read from config file, some kine from commandline.
It's not convient if I want to specify a single configuration with .etc. commandline, if it's had been using config file.

Package config provides an unify method for different config sources.
*/
package config
