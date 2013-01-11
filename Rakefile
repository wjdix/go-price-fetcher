def abs_path(path)
	File.expand_path(File.join(current_path, path))
end

def current_path
	File.expand_path(File.dirname(__FILE__))
end

def go_command(command)
	sh "env GOPATH=#{current_path} #{command}"
end

desc "install dependencies"
task :deps do
	dependencies = ["code.google.com/p/go-html-transform/html/transform", "github.com/moovweb/gokogiri"]
	dependencies.each do |dependency|
		go_command("go get #{dependency}")
	end
end

desc "clean compiled files and binaries"
task :clean do
	dirs = %w[bin pkg]
	dirs.each do |dir|
		rm_rf abs_path(dir)
	end
end

desc "format code"
task :fmt do
	go_command("go fmt wdix/getev/...")
end

desc "build the project"
task :build => [:deps, :clean, :fmt] do
	go_command("go install wdix/getev/...")
end

task :default => :build
