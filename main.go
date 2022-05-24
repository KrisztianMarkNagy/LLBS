package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type llbs_config struct {
	Compiler            string   `json:"compiler"`
	Linker              string   `json:"linker"`
	Project_name        string   `json:"project_name"`
	Object_files_dir    string   `json:"object_files_dir"`
	Executable_dir      string   `json:"executable_dir"`
	Source_dir          string   `json:"source_dir"`
	Local_headers_dir   string   `json:"local_headers_dir"`
	Local_mods_dir      string   `json:"local_mods_dir"`
	Local_libs_dirs     []string `json:"local_libs_dirs"`
	Extern_headers_dirs []string `json:"extern_headers_dirs"`
	Extern_mods_dirs    []string `json:"extern_mods_dirs"`
	Extern_libs_dirs    []string `json:"extern_libs_dirs"`
	Excluded_dirs       []string `json:"excluded_dirs"`
	Liblinks            []string `json:"liblinks"`
	Warnings            []string `json:"warnings"`
	Flags               []string `json:"flags"`
	Other               []string `json:"other"`
}

const (
	idle_stage = iota
	comp_stage = iota
	link_stage = iota
)

var (
	current_stage = idle_stage
	config        llbs_config
	build_file    string
)

func (c *llbs_config) from_file(fpath string) error {
	file_content_in_bytes, err := os.ReadFile(fpath)
	if err != nil {
		log.SetPrefix("\033[33m[Warning]\033[0m ")
		log.Printf("Error occured while trying to read from file `%s`, from within the `llbs_config.from_file()` method.\n", fpath)
		log.SetPrefix("")
		return err
	}
	reg, err := regexp.Compile("[\t ]{0,}//.+")
	if err != nil {
		return err
	}
	file_content_in_bytes = []byte(reg.ReplaceAllLiteralString(string(file_content_in_bytes), "\r"))
	// fmt.Println("LLBS Settings:\n" + string(file_content_in_bytes))

	if err = json.Unmarshal(file_content_in_bytes, c); err != nil {
		log.SetPrefix("\033[33m[Warning]\033[0m ")
		log.Println("Error occured inside the json unmarshalling condition, from within the `llbs_config.from_file()` method.")
		log.SetPrefix("")
		return err
	}

	c.Object_files_dir, err = filepath.Abs(filepath.Clean(c.Object_files_dir))
	if err != nil {
		return err
	}

	c.Executable_dir, err = filepath.Abs(filepath.Clean(c.Executable_dir))
	if err != nil {
		return err
	}

	c.Source_dir, err = filepath.Abs(filepath.Clean(c.Source_dir))
	if err != nil {
		return err
	}

	c.Local_headers_dir, err = filepath.Abs(filepath.Clean(c.Local_headers_dir))
	if err != nil {
		return err
	}

	c.Local_mods_dir, err = filepath.Abs(filepath.Clean(c.Local_mods_dir))
	if err != nil {
		return err
	}

	for idx := range c.Local_libs_dirs {
		if !fs.ValidPath(c.Local_libs_dirs[idx]) {
			return fmt.Errorf("path `%s` is invalid", c.Local_libs_dirs[idx])
		}
		c.Local_libs_dirs[idx], err = filepath.Abs(filepath.Clean(c.Local_libs_dirs[idx]))
		if err != nil {
			return err
		}
	}

	for idx := range c.Extern_headers_dirs {
		if !fs.ValidPath(c.Extern_headers_dirs[idx]) {
			return fmt.Errorf("path `%s` is invalid", c.Extern_headers_dirs[idx])
		}
		c.Extern_headers_dirs[idx], err = filepath.Abs(filepath.Clean(c.Extern_headers_dirs[idx]))
		if err != nil {
			return err
		}
	}

	for idx := range c.Extern_mods_dirs {
		if !fs.ValidPath(c.Extern_mods_dirs[idx]) {
			return fmt.Errorf("path `%s` is invalid", c.Extern_mods_dirs[idx])
		}
		c.Extern_mods_dirs[idx], err = filepath.Abs(filepath.Clean(c.Extern_mods_dirs[idx]))
		if err != nil {
			return err
		}
	}

	for idx := range c.Extern_libs_dirs {
		if !fs.ValidPath(c.Extern_libs_dirs[idx]) {
			return fmt.Errorf("path `%s` is invalid", c.Extern_libs_dirs[idx])
		}
		c.Extern_libs_dirs[idx], err = filepath.Abs(filepath.Clean(c.Extern_libs_dirs[idx]))
		if err != nil {
			return err
		}
	}

	for idx := range c.Excluded_dirs {
		if !fs.ValidPath(c.Excluded_dirs[idx]) {
			return fmt.Errorf("path `%s` is invalid", c.Excluded_dirs[idx])
		}
		c.Excluded_dirs[idx], err = filepath.Abs(filepath.Clean(c.Excluded_dirs[idx]))
		if err != nil {
			return err
		}
	}

	fmt.Println("\033[34m[Process]\033[0m Starting compilation process...")
	return c.compilation_stage()
}

func (c llbs_config) compilation_stage() error {
	var (
		external_libraries []string
		external_modules   []string
		local_libraries    []string
		local_modules      []string
		source_files       []string
		err                error
		arguments          []string = []string{"-c"}
	)

	if len(c.Other) != 0 {
		arguments = append(arguments, c.Other...)
	}
	if len(c.Flags) != 0 {
		arguments = append(arguments, c.Flags...)
	}
	if len(c.Warnings) != 0 {
		arguments = append(arguments, c.Warnings...)
	}
	if len(c.Liblinks) != 0 {
		arguments = append(arguments, c.Liblinks...)
	}
	if len(c.Extern_headers_dirs) != 0 {
		arguments = append(arguments, "-I="+strings.Join(c.Extern_headers_dirs, " -I="))
	}
	if len(c.Local_headers_dir) != 0 {
		arguments = append(arguments, "-I="+c.Local_headers_dir)
	}

	if err = c.add_to_list(&external_libraries, c.Extern_libs_dirs); err != nil {
		return err
	}

	if err = c.add_to_list(&external_modules, c.Extern_mods_dirs); err != nil {
		return err
	}

	if err = c.add_to_list(&local_libraries, c.Local_libs_dirs); err != nil {
		return err
	}

	if err = c.add_to_list(&local_modules, c.Local_mods_dir); err != nil {
		return err
	}

	if err = c.add_to_list(&source_files, c.Source_dir); err != nil {
		return err
	}

	if err = c.compile_cfile(&arguments, external_libraries); err != nil {
		fmt.Println("\033[31m[Error]\033[0m Something occured while compiling `external_libraries`")
		return err
	}
	if err = c.compile_cfile(&arguments, external_modules); err != nil {
		fmt.Println("\033[31m[Error]\033[0m Something occured while compiling `external_modules`")
		return err
	}
	if err = c.compile_cfile(&arguments, local_libraries); err != nil {
		fmt.Println("\033[31m[Error]\033[0m Something occured while compiling `local_libraries`")
		return err
	}
	if err = c.compile_cfile(&arguments, local_modules); err != nil {
		fmt.Println("\033[31m[Error]\033[0m Something occured while compiling `local_modules`")
		return err
	}
	if err = c.compile_cfile(&arguments, source_files); err != nil {
		fmt.Println("\033[31m[Error]\033[0m Something occured while compiling `source_files`")
		return err
	}
	// fmt.Printf("\033[34m[Info]\033[0m Common Arguments: %q\n", arguments)

	fmt.Println("\n\033[34m[Process]\033[0m Starting linking process...")
	return c.linking_stage()
}

func (c llbs_config) linking_stage() error {
	arguments := []string{fmt.Sprintf("-o%s/%s", c.Executable_dir, c.Project_name)}

	dircnt, err := ioutil.ReadDir(c.Object_files_dir)
	if err != nil {
		return err
	}
	for _, vol := range dircnt {
		if vol.IsDir() {
			continue
		}
		if !strings.HasSuffix(vol.Name(), ".o") {
			continue
		}
		arguments = append(arguments, fmt.Sprintf("%s/%s", c.Object_files_dir, vol.Name()))
	}

	process := exec.Command(c.Linker, arguments...)

	fmt.Println("\033[34m[Info]\033[0m\t Command:", process.Path, process.Args)

	output, err := process.CombinedOutput()
	if err != nil {
		return err
	}

	exit_code := process.ProcessState.ExitCode()
	if len(string(output)) != 0 {
		if exit_code == 0 {
			fmt.Printf("\n\033[32m[Output]:\033[0m\n%s\n\033[32m[Status]:\033[0m %d\n\n", string(output), exit_code)
		} else {
			fmt.Printf("\n\033[31m[Output]:\033[0m\n%s\n\033[31m[Status]:\033[0m %d\n\n", string(output), exit_code)
		}
	} else {
		if exit_code == 0 {
			fmt.Printf("\n\033[32m[Status] (Exit Code):\033[0m %d\n\n", exit_code)
		} else {
			fmt.Printf("\n\033[31m[Status] (Exit Code):\033[0m %d\n\n", exit_code)
		}
	}
	return nil
}

// * For the compilation stage:
func (c llbs_config) add_to_list(lst *[]string, from interface{}) error {
	var err error = nil
	switch from := from.(type) {
	case []string:
		for idx := range from {
			if err = filepath.WalkDir(from[idx], func(path string, d fs.DirEntry, err error) error {
				var _err error
				for _, excl_path := range c.Excluded_dirs {
					if !filepath.IsAbs(excl_path) {
						excl_path, _err = filepath.Abs(filepath.Clean(excl_path))
						if _err != nil {
							return _err
						}
					}
					if excl_path == path {
						return nil
					}
				}
				if filepath.Ext(path) == ".c" {
					*lst = append(*lst, path)
				}
				return nil
			}); err != nil {
				return err
			}
		}

	case string:
		if err = filepath.WalkDir(from, func(path string, d fs.DirEntry, err error) error {
			var _err error
			for _, excl_path := range c.Excluded_dirs {
				if !filepath.IsAbs(excl_path) {
					excl_path, _err = filepath.Abs(filepath.Clean(excl_path))
					if _err != nil {
						return _err
					}
				}
				if excl_path == path {
					return nil
				}
			}
			if filepath.Ext(path) == ".c" {
				*lst = append(*lst, path)
			}
			return nil
		}); err != nil {
			return err
		}

	default:
		return fmt.Errorf("invalid type of `from` parameter inside `llbs_config.add_to_list`")
	}

	return nil
}

func (c llbs_config) compile_cfile(args *[]string, over []string) error {
	for _, cfile_path := range over {
		_args := *args
		fmt.Println("\033[34m[Info]\033[0m Compiling:", cfile_path)
		_, cfile := filepath.Split(cfile_path)
		_args = append(_args, "-o"+filepath.Join(c.Object_files_dir, cfile[:len(cfile)-2]+".o"), cfile_path)
		process := exec.Command(c.Compiler, _args...)
		fmt.Println("\033[34m[Info]\033[0m\t Command:", process.Path, process.Args)

		output, err := process.CombinedOutput()
		if err != nil {
			return err
		}

		exit_code := process.ProcessState.ExitCode()
		if len(string(output)) != 0 {
			if exit_code == 0 {
				fmt.Printf("\n\033[32m[Output]:\033[0m\n%s\n\033[32m[Status]:\033[0m %d\n\n", string(output), exit_code)
			} else {
				fmt.Printf("\n\033[31m[Output]:\033[0m\n%s\n\033[31m[Status]:\033[0m %d\n\n", string(output), exit_code)
			}
		} else {
			if exit_code == 0 {
				fmt.Printf("\n\033[32m[Status] (Exit Code):\033[0m %d\n\n", exit_code)
			} else {
				fmt.Printf("\n\033[31m[Status] (Exit Code):\033[0m %d\n\n", exit_code)
			}
		}
	}
	return nil
}

// * For the compilation stage:
// ...

// * Main thread, main process:
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmsgprefix)
	flag.StringVar(&build_file, "f", "./build.llbs.jsonc", "Build file input")
	flag.Parse()
}

// TODO: Develop a build-system for C/C++/Assembly?
func main() {
	fmt.Print("\n\t\033[32m!START!\033[0m\n\n")

	if err := config.from_file(build_file); err != nil {
		log.SetPrefix("\033[31m[Error]\033[0m ")
		log.Fatalln(err.Error())
		return
	}

	fmt.Println("\033[32m[Success]\033[0m The executable binary can be found at:", fmt.Sprintf("%s/%s", config.Executable_dir, config.Project_name))

	fmt.Print("\n\t\033[33m!DONE!\033[0m\n\n")
}
