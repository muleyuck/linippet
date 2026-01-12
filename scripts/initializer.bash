linippet_apply() {
    local snippet="$(linippet)"

    if [[ -z $snippet ]]; then
        return 1
    fi

    eval "$snippet"

    return 0
}
alias lip=linippet_apply

export LINIPPET_TRIGGER_BIND_KEY=${LINIPPET_TRIGGER_BIND_KEY}

# READLINE is supported at version which is 4 or later
if [[ -n $LINIPPET_TRIGGER_BIND_KEY && -n "${BASH_VERSINFO[0]}" && "${BASH_VERSINFO[0]}" -ge 4 ]]; then
    linippet_triggered() {
        local snippet="$(linippet)"

        if [[ -z $snippet ]]; then
            return 1
        fi

        READLINE_LINE="${snippet}"
        READLINE_POINT=${#READLINE_LINE}

        return 0
    }
    bind -x "\"$LINIPPET_TRIGGER_BIND_KEY\": linippet_triggered"

fi
