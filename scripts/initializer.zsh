linippet_apply() {
    local snippet="$(linippet)"

    if [[ -z $snippet ]]; then
        return 1
    fi

    # Only One-liner
    if [[ $snippet == *"\n"* || $snippet == *"\r"* ]]; then
        echo "linippet is supported only one-liner snippet"
        return 1
    fi

    eval "$snippet"

    return 0
}
alias lip=linippet_apply

export LINIPPET_TRIGGER_BIND_KEY=${LINIPPET_TRIGGER_BIND_KEY}

if [[ -n $LINIPPET_TRIGGER_BIND_KEY ]]; then
    linippet_triggered() {
        local snippet="$(linippet)"

        if [[ -z $snippet ]]; then
            zle accept-line
            return 1
        fi

        # Only One-liner
        if [[ $snippet == *"\n"* || $snippet == *"\r"* ]]; then
            echo "linippet is supported only one-liner snippet"
            zle accept-line
            return 1
        fi

        LBUFFER="${snippet}"
        CURSOR=$#BUFFER
        zle reset-prompt

        return 0
    }

    zle -N linippet_triggered
    bindkey ${LINIPPET_TRIGGER_BIND_KEY} linippet_triggered
fi
