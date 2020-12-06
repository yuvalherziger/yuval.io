package main

import (
    "fmt"
    "github.com/jedib0t/go-pretty/list"
    "github.com/jedib0t/go-pretty/table"
    "github.com/jedib0t/go-pretty/text"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "log"
    "os"
    "sort"
    "strings"
)

var (
    configFile string
    logLevel   string
    rootCmd    *cobra.Command
)

const (
    defaultConfigPath     string = "/etc/config/"
    defaultConfigFilename string = "config"
    banner                string = `
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃   THE                          ┃
┃   ╭━━╮╭━━┳━━━╮  ╭━━━┳╮╱╱╭━━╮   ┃
┃   ┃╭╮┃╰┫┣┫╭━╮┃  ┃╭━╮┃┃╱╱╰┫┣╯   ┃
┃   ┃╰╯╰╮┃┃┃┃╱┃┃  ┃┃╱╰┫┃╱╱╱┃┃    ┃
┃   ┃╭━╮┃┃┃┃┃╱┃┃  ┃┃╱╭┫┃╱╭╮┃┃    ┃
┃   ┃╰━╯┣┫┣┫╰━╯┃  ┃╰━╯┃╰━╯┣┫┣╮   ┃ 
┃   ╰━━━┻━━┻━━━╯  ╰━━━┻━━━┻━━╯   ┃
┃              BY YUVAL HERZIGER ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛

Your interface to my CV; read the manual below to find out more:
`
)

type Profile struct {
    Name string `mapstructure:"name"`
    Url  string `mapstructure:"url"`
}

type Skill struct {
    Name        string `mapstructure:"name"`
    Proficiency int    `mapstructure:"proficiency"`
}

type Role struct {
    Name             string   `mapstructure:"name"`
    StartDate        string   `mapstructure:"startDate"`
    EndDate          string   `mapstructure:"endDate"`
    Responsibilities []string `mapstructure:"responsibilities"`
}

type Company struct {
    CompanyName    string `mapstructure:"companyName"`
    CompanyWebsite string `mapstructure:"companyWebsite"`
    Location       string `mapstructure:"location"`
    Roles          []Role `mapstructure:"roles"`
}

type Certification struct {
    Name        string `mapstructure:"name"`
    Credentials string `mapstructure:"credentials"`
    Expires     string `mapstructure:"expires"`
}

type Education struct {
    Institute      string `mapstructure:"institute"`
    Degree         string `mapstructure:"degree"`
    Major          string `mapstructure:"major"`
    GraduationDate string `mapstructure:"graduationDate"`
    GPA            string `mapstructure:"GPA"`
}

func ParseConfig() error {
    return nil
}

func initConfig() {
    if configFile == "" {
        viper.SetConfigType("yaml")
        viper.AddConfigPath(defaultConfigPath)
        viper.SetConfigName(defaultConfigFilename)
    } else {
        viper.SetConfigFile(configFile)
    }

    err := viper.ReadInConfig()
    if err != nil {
        log.Fatalf("Failed to read config file: %v", err)
    }
    err = ParseConfig()
    if err != nil {
        log.Fatalf("Failed to parse the service configurations: %v", err)
    }
}

func getUnitSize(vWidth int, cols int, divider int) int {
    padding := (cols-1)*3 + 4
    return (vWidth - padding) / divider
}

func experience(cmd *cobra.Command, args []string) {
    nRoles, err := cmd.PersistentFlags().GetInt("roles")
    if err != nil {
        panic(err)
    }
    vWidth, err := cmd.PersistentFlags().GetInt("view-width")
    if err != nil {
        panic(err)
    }
    pretty, err := cmd.PersistentFlags().GetBool("pretty")
    if err != nil {
        panic(err)
    }
    t := table.NewWriter()
    if pretty {
        t.SetStyle(table.StyleColoredBlackOnCyanWhite)
    } else {
        t.SetStyle(table.StyleLight)
    }
    t.SetOutputMirror(os.Stdout)
    t.SetAutoIndex(true)
    t.Style().Options.DrawBorder = true
    t.Style().Options.SeparateRows = true
    unit := getUnitSize(vWidth, 4, 6)
    t.SetColumnConfigs([]table.ColumnConfig{
        {Name: "Company", WidthMax: unit},
        {Name: "Role", WidthMax: unit},
        {Name: "Responsibilities", WidthMax: unit * 4, Align: text.AlignLeft},
    })
    t.AppendHeader(table.Row{"Company", "Role", "Responsibilities"})
    var companies []Company
    err = viper.UnmarshalKey("experience", &companies)
    if err != nil {
        panic(err)
    }
    var prevC Company
    tRoles := 0
    for _, c := range companies {
        for _, r := range c.Roles {
            cName := ""
            cLoc := ""
            cWeb := ""
            if c.CompanyName != prevC.CompanyName {
                cName = c.CompanyName
                cLoc = c.Location
                cWeb = c.CompanyWebsite
            }
            l := list.NewWriter()
            l.SetStyle(list.StyleBulletTriangle)
            for _, rs := range r.Responsibilities {
                l.AppendItem(rs)
            }

            t.AppendRows([]table.Row{
                {
                    fmt.Sprintf("%s\n%s\n%s", cName, cLoc, cWeb),
                    fmt.Sprintf("%s\n%s - %s", r.Name, r.StartDate, r.EndDate),
                    fmt.Sprintf("%s\n", l.Render()),
                },
            })
            prevC = c
            if nRoles > 0 && (tRoles+1) >= nRoles {
                t.Render()
                return
            }
            tRoles++
        }
    }
    t.Render()
}

func about(cmd *cobra.Command, args []string) {
    pretty, err := cmd.PersistentFlags().GetBool("pretty")
    if err != nil {
        panic(err)
    }
    vWidth, err := cmd.PersistentFlags().GetInt("view-width")
    if err != nil {
        panic(err)
    }
    t := table.NewWriter()
    if pretty {
        t.SetStyle(table.StyleColoredBlackOnCyanWhite)
    } else {
        t.SetStyle(table.StyleLight)
    }
    t.SetOutputMirror(os.Stdout)
    t.Style().Options.DrawBorder = true
    t.Style().Options.SeparateRows = true
    unit := getUnitSize(vWidth, 2, 5)
    t.SetColumnConfigs([]table.ColumnConfig{
        {Number: 0, WidthMax: unit},
        {Number: 1, WidthMax: unit * 4},
    })
    t.AppendRows([]table.Row{{
        "Name", fmt.Sprintf("%s (%s) %s", viper.GetString("about.firstName"), viper.GetString("about.nickname"), viper.GetString("about.lastName")),
    }})
    t.AppendRows([]table.Row{{
        "About", viper.GetString("about.intro"),
    }})
    t.AppendRows([]table.Row{{
        "Contact", viper.GetString("about.contact.emailAddress"),
    }})
    pList := list.NewWriter()
    pList.SetStyle(list.StyleBulletTriangle)
    var profiles []Profile
    err = viper.UnmarshalKey("about.profiles", &profiles)
    if err != nil {
        panic(err)
    }
    for _, p := range profiles {
        pList.AppendItem(fmt.Sprintf("%s: %s", p.Name, p.Url))
    }
    t.AppendRows([]table.Row{{
        "Online Profiles", pList.Render(),
    }})

    lList := list.NewWriter()
    lList.SetStyle(list.StyleBulletTriangle)
    for _, lang := range viper.GetStringSlice("about.languages") {
        lList.AppendItem(lang)
    }
    t.AppendRows([]table.Row{{
        "Languages", lList.Render(),
    }})
    t.Render()
}

func picture(cmd *cobra.Command, args []string) {
    fmt.Printf(viper.GetString("picture"))
}

func buildSkillProgress(p int) string {
    return fmt.Sprintf(
        "%s%s %d/%d",
        strings.Repeat("█", p*2),
        strings.Repeat("░", (10-p)*2),
        p,
        10,
    )
}

func certifications(cmd *cobra.Command, args []string) {
    pretty, err := cmd.PersistentFlags().GetBool("pretty")
    if err != nil {
        panic(err)
    }
    vWidth, err := cmd.PersistentFlags().GetInt("view-width")
    if err != nil {
        panic(err)
    }
    var certifications []Certification
    err = viper.UnmarshalKey("certifications", &certifications)
    if err != nil {
        panic(err)
    }

    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
    t.SetAutoIndex(true)
    t.Style().Options.DrawBorder = true
    t.Style().Options.SeparateRows = true
    if pretty {
        t.SetStyle(table.StyleColoredBlackOnCyanWhite)
    } else {
        t.SetStyle(table.StyleLight)
    }
    // 1, 4 , 5
    unit := getUnitSize(vWidth, 3, 10)
    t.SetColumnConfigs([]table.ColumnConfig{
        {Name: "Certification", WidthMax: unit * 4},
        {Name: "Expires", WidthMax: unit},
        {Name: "Credentials", WidthMax: unit * 5},
    })
    t.AppendHeader(table.Row{"Certification", "Expires", "Credentials"})
    for _, c := range certifications {
        t.AppendRows([]table.Row{{
            c.Name,
            c.Expires,
            c.Credentials,
        }})
    }
    t.Render()
}

func education(cmd *cobra.Command, args []string) {
    pretty, err := cmd.PersistentFlags().GetBool("pretty")
    if err != nil {
        panic(err)
    }
    vWidth, err := cmd.PersistentFlags().GetInt("view-width")
    if err != nil {
        panic(err)
    }
    var education []Education
    err = viper.UnmarshalKey("education", &education)
    if err != nil {
        panic(err)
    }

    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
    t.SetAutoIndex(true)
    t.Style().Options.DrawBorder = true
    t.Style().Options.SeparateRows = true
    if pretty {
        t.SetStyle(table.StyleColoredBlackOnCyanWhite)
    } else {
        t.SetStyle(table.StyleLight)
    }

    unit := getUnitSize(vWidth, 5, 1)
    t.SetColumnConfigs([]table.ColumnConfig{
        {Name: "Institute", WidthMax: unit},
        {Name: "Degree", WidthMax: unit},
        {Name: "Major", WidthMax: unit},
        {Name: "Graduation", WidthMax: unit},
        {Name: "GPA", WidthMax: unit},
    })
    t.AppendHeader(table.Row{"Institute", "Degree", "Major", "Graduation", "GPA"})
    for _, deg := range education {
        t.AppendRows([]table.Row{{
            deg.Institute,
            deg.Degree,
            deg.Major,
            deg.GraduationDate,
            deg.GPA,
        }})
    }
    t.Render()
}

func openSource(cmd *cobra.Command, args []string) {

}

func openProfile(cmd *cobra.Command, args []string) {
    var profiles []Profile
    _ = viper.UnmarshalKey("about.profiles", &profiles)
    profileName := args[0]
    for _, p := range profiles {
        if strings.ToLower(profileName) == strings.ToLower(p.Name) {
            fmt.Printf("Opening %s", p.Url)
            return
        }
    }
    fmt.Printf("Profile not found")
    return
}

func skills(cmd *cobra.Command, args []string) {
    pretty, err := cmd.PersistentFlags().GetBool("pretty")
    if err != nil {
        panic(err)
    }
    vWidth, err := cmd.PersistentFlags().GetInt("view-width")
    if err != nil {
        panic(err)
    }
    var skills []Skill
    err = viper.UnmarshalKey("skills", &skills)
    if err != nil {
        panic(err)
    }
    sort.SliceStable(skills, func(i, j int) bool {
        return skills[i].Proficiency > skills[j].Proficiency
    })
    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
    t.SetAutoIndex(true)
    t.Style().Options.DrawBorder = true
    t.Style().Options.SeparateRows = true
    if pretty {
        t.SetStyle(table.StyleColoredBlackOnCyanWhite)
    } else {
        t.SetStyle(table.StyleLight)
    }
    unit := getUnitSize(vWidth, 2, 3)
    t.SetColumnConfigs([]table.ColumnConfig{
        {Name: "Skill", WidthMax: unit * 2},
        {Name: "Proficiency", WidthMax: unit},
    })
    t.AppendHeader(table.Row{"Skill", "Proficiency"})
    for _, s := range skills {
        t.AppendRows([]table.Row{{
            s.Name,
            buildSkillProgress(s.Proficiency),
        }})
    }
    t.Render()

}

func init() {
    cobra.OnInitialize(initConfig)
    rootCmd = &cobra.Command{
        Use:        "bio",
        Short:      fmt.Sprintf("%s\n\n%s", banner, viper.GetString("intro")),
        SuggestFor: []string{"bio"},
        Version:    "v1beta1",
    }

    fs := rootCmd.Flags()
    fs.Int("view-width", 256, "View width")
    _ = fs.MarkHidden("view-width")

    expCmd := &cobra.Command{
        Use:     "experience",
        Aliases: []string{"e"},
        Short:   "Show my work experience",
        Run:     experience,
    }
    expFlags := expCmd.PersistentFlags()
    expFlags.Int("roles", -1, "The # of roles to retrieve, in descending chronological order")
    expFlags.Bool("pretty", false, "Return a colorful table")
    expFlags.Int("view-width", 256, "View width")
    _ = expFlags.MarkHidden("view-width")

    aboutCmd := &cobra.Command{
        Use:     "about",
        Aliases: []string{"a"},
        Short:   "Show information about me",
        Run:     about,
    }
    aboutFlags := aboutCmd.PersistentFlags()
    aboutFlags.Bool("pretty", false, "Return a colorful table")
    aboutFlags.Int("view-width", 256, "View width")
    _ = aboutFlags.MarkHidden("view-width")

    picCmd := &cobra.Command{
        Use:     "picture",
        Aliases: []string{"p", "pic"},
        Short:   "Show my picture",
        Run:     picture,
    }
    picFlags := picCmd.PersistentFlags()
    picFlags.Int("view-width", 256, "View width")
    _ = picFlags.MarkHidden("view-width")

    sklCmd := &cobra.Command{
        Use:     "skills",
        Aliases: []string{"s"},
        Short:   "Show my self-proclaimed skills ;-)",
        Run:     skills,
    }
    sklFlags := sklCmd.PersistentFlags()
    sklFlags.Bool("pretty", false, "Return a colorful table")
    sklFlags.Int("view-width", 256, "View width")
    _ = sklFlags.MarkHidden("view-width")

    crtCmd := &cobra.Command{
        Use:     "certifications",
        Aliases: []string{"c"},
        Short:   "Show my certifications",
        Run:     certifications,
    }
    crtFlags := crtCmd.PersistentFlags()
    crtFlags.Bool("pretty", false, "Return a colorful table")
    crtFlags.Int("view-width", 256, "View width")
    _ = crtFlags.MarkHidden("view-width")

    eduCmd := &cobra.Command{
        Use:     "education",
        Aliases: []string{"ed"},
        Short:   "Show my education",
        Run:     education,
    }
    eduFlags := eduCmd.PersistentFlags()
    eduFlags.Bool("pretty", false, "Return a colorful table")
    eduFlags.Int("view-width", 256, "View width")
    _ = eduFlags.MarkHidden("view-width")
    initConfig()
    var profiles []Profile
    _ = viper.UnmarshalKey("about.profiles", &profiles)
    var pNames []string
    for _, p := range profiles {
        pNames = append(pNames, strings.ToLower(p.Name))
    }

    openCmd := &cobra.Command{
        Use:       "open",
        Aliases:   []string{"o"},
        Short:     "Open a profile page in a new tab",
        Run:       openProfile,
        ValidArgs: pNames,
        Args:      cobra.ExactValidArgs(1),
        Example:   fmt.Sprintf("bio open [ %s ]", strings.Join(pNames[:], " | ")),
    }
    openFlags := openCmd.PersistentFlags()
    openFlags.Int("view-width", 256, "View width")
    _ = openFlags.MarkHidden("view-width")

    rootCmd.AddCommand(expCmd)
    rootCmd.AddCommand(picCmd)
    rootCmd.AddCommand(crtCmd)
    rootCmd.AddCommand(aboutCmd)
    rootCmd.AddCommand(sklCmd)
    rootCmd.AddCommand(eduCmd)
    rootCmd.AddCommand(openCmd)
}

func main() {
    err := rootCmd.Execute()
    if err != nil {
        log.Fatalf("Failed to run: %v", err)
    }
}
