function Setup()
    -- SetEnv(EnvName string, Value string)
    -- SetVar(VerName string, Value string)
    -- SetExt(ExtName string, Pattern string)
    -- SetName(ExtName string)

    SetEnv("temp", "Temp")
    SetEnv("exec-rule", "async")

    SetVar("Whisper", '"{@edr}/Tools/ToWhisperNet/ToWhisperNet.exe"')
    SetVar("Whisper-l", 0.9)
    SetVar("Renamer", '"{@edr}/Tools/Renamer/renamer.exe"')
    SetVar("NewLine", '"{@edr}/Tools/auto-NewLine/auto-NewLine.exe"')
    SetVar("SoundEffect", '"{@edr}/Tools/sox-14.4.2/sox.exe"')
    SetVar("Numbering", '"{@edr}/Tools/Numbering/Numbering.exe"')

    SetExt("vic", "\\.wav$")
    SetExt("txt", "(\\.text$)|(\\.txt$)")

    SetName("vic")
end

function Process()
    -- Match(Pattern string, TextFilePath ?string, Encode ?string["utf-8" | "shift-jis"])
    -- Execute(Command: string)
    -- Wait()
    -- SetExt(ExtName string, Pattern string)
    -- ClearExt()

    Execute('{@Renamer} -t "{@src}/{@txt}" "{@src}/{@vic}"')
    Execute('{@NewLine} -min 10 -max 30 -t "{@src}/{@txt}"')

    if Match('^.*?_[”"]') then
        Execute('{@SoundEffect} "{@src}/{@vic}" "{@dst}/{@vic}" vol 0.5')
        Execute('{@Whisper} -l {@whisper-l} -o "{@dst}/{@vic}" "{@src}/{@vic}"')
    end

    if Match('^.*?_[（(]') then
        Execute('{@SoundEffect} "{@src}/{@vic}" "{@dst}/{@vic}" echo 1 0.6 100 0.25')
    end

    if Match('^.*?_[＃#]') then
        Execute('{@SoundEffect} "{@src}/{@vic}" "{@dst}/{@vic}" bandpass 2306 250h vol 10')
    end

    if Match('^.*?_[｜|]') then
        Execute('{@SoundEffect} "{@src}/{@vic}" "{@dst}/{@vic}" lowpass 1700 10000h vol 1.6')
    end

    Wait()
    Execute('{@Numbering} -d "{@out}" "{@src}/{@vic}" "{@src}/{@txt}"')
end

function Notify()
    -- Match(Pattern string, TextFilePath ?string, Encode ?string["utf-8" | "shift-jis"])
    -- Execute(Command: string)

    -- Execute('xxx.exe "{@vic}" "{@txt}"')
end
