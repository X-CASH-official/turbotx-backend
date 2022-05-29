# Color print variables
COLOR_PRINT_RED ?= "\033[1;31m"
COLOR_PRINT_GREEN ?= "\033[1;32m"
END_COLOR_PRINT ?= "\033[0m"

# The binary name
TARGET_BINARY ?= turbotxmodule

# The build directory, where the target binary will be stored
BUILD_DIR ?= ./

# The source directory, where the 
SRC_DIRS ?= ./

# Linker flags
LDFLAGS ?= 

# Set the compiler and linker flags for the different options
debug: 
release: LDFLAGS += -s -w

# Set the debug and release rules to phony, since they just change the compiler flags variable and dont create any files. Set the clean rule to phony since we want it to run the command everytime we run it
.PHONY: debug release clean

# Set the options to do the same thing, since the only difference is the compiler flags
debug: $(BUILD_DIR)$(TARGET_BINARY)
release: $(BUILD_DIR)$(TARGET_BINARY)

# Link all of the objects files
$(BUILD_DIR)$(TARGET_BINARY):
	@go build -ldflags="$(LDFLAGS)"
	@echo "\n" $(COLOR_PRINT_GREEN)$(TARGET_BINARY) "Has Been Built Successfully"$(END_COLOR_PRINT)

# Remove the build directory
clean:
	@$(RM) $(BUILD_DIR)$(TARGET_BINARY)
	@echo $(COLOR_PRINT_RED) "Removed" $(BUILD_DIR)$(TARGET_BINARY)$(END_COLOR_PRINT)
