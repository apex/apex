#compdef apex

_message_next_arg()
{
    argcount=0
    for word in "${words[@][2,-1]}"
    do
        if [[ $word != -* ]] ; then
            ((argcount++))
        fi
    done
    if [[ $argcount -le ${#myargs[@]} ]] ; then
        _message -r $myargs[$argcount]
        if [[ $myargs[$argcount] =~ ".*file.*" || $myargs[$argcount] =~ ".*path.*" ]] ; then
            _files
        fi
    fi
}

_apex ()
{
    local context state state_descr line
    typeset -A opt_args

    _arguments -C \
        ':command:->command' \
		'(-h)-h[Output help information]' \
		'(--help)--help[Output help information]' \
		'(-h)-h[Output help information]' \
		'(--help)--help[Output help information]' \
		'(-V)-V[Output version]' \
		'(--version)--version[Output version]' \
        '*::options:->options'

    case $state in
        (command)
            local -a subcommands
            subcommands=(
				'rollback'
				'logs'
				'invoke'
				'deploy'
				'list'
				'build'
				'delete'
            )
            _values 'apex' $subcommands
        ;;

        (options)
            case $line[1] in
                rollback)
                    _apex-rollback
                ;;
                logs)
                    _apex-logs
                ;;
                invoke)
                    _apex-invoke
                ;;
                deploy)
                    _apex-deploy
                ;;
                list)
                    _apex-list
                ;;
                build)
                    _apex-build
                ;;
                delete)
                    _apex-delete
                ;;
            esac
        ;;
    esac

}

_apex-rollback ()
{
    local context state state_descr line
    typeset -A opt_args

    if [[ $words[$CURRENT] == -* ]] ; then
        _arguments -C \
        ':command:->command' \

    else
        myargs=('<name>' '<version>')
        _message_next_arg
    fi
}

_apex-logs ()
{
    local context state state_descr line
    typeset -A opt_args

    if [[ $words[$CURRENT] == -* ]] ; then
        _arguments -C \
        ':command:->command' \
		'(-F=-)-F=-' \
		'(--filter=-)--filter=-' \

    else
        myargs=('<name>')
        _message_next_arg
    fi
}

_apex-invoke ()
{
    local context state state_descr line
    typeset -A opt_args

    if [[ $words[$CURRENT] == -* ]] ; then
        _arguments -C \
        ':command:->command' \
		'(-a)-a[Async invocation]' \
		'(--async)--async[Async invocation]' \
		'(-v)-v[Output verbose logs]' \
		'(--verbose)--verbose[Output verbose logs]' \

    else
        myargs=('<name>')
        _message_next_arg
    fi
}

_apex-deploy ()
{
    local context state state_descr line
    typeset -A opt_args

    if [[ $words[$CURRENT] == -* ]] ; then
        _arguments -C \
        ':command:->command' \

    else
        myargs=('<name>')
        _message_next_arg
    fi
}

_apex-list ()
{
    local context state state_descr line
    typeset -A opt_args

    _arguments -C \
        ':command:->command' \
        
}

_apex-build ()
{
    local context state state_descr line
    typeset -A opt_args

    if [[ $words[$CURRENT] == -* ]] ; then
        _arguments -C \
        ':command:->command' \

    else
        myargs=('<name>')
        _message_next_arg
    fi
}

_apex-delete ()
{
    local context state state_descr line
    typeset -A opt_args

    if [[ $words[$CURRENT] == -* ]] ; then
        _arguments -C \
        ':command:->command' \

    else
        myargs=('<name>')
        _message_next_arg
    fi
}


_apex "$@"