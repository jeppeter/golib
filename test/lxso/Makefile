

include ../../baseop.mak

SYSTEM_NAME:=$(shell uname)


TARGET  = libsomain

HEADERS = somain.h
SOURCES = somain.c
OBJECTS = $(patsubst %.c,%.o,$(SOURCES))

CFLAGS  = $(COMMON_CFLAGS)
CFLAGS += -Wall
CFLAGS += -fPIC 
LDFLAGS = 

ifeq ($(SYSTEM_NAME), Linux)
LDFLAGS += -shared -Wl,-soname,$(TARGET).so
endif

ifeq ($(SYSTEM_NAME), Darwin)
LDFLAGS += -dynamiclib -install_name @rpath/$(TARGET).so
endif

all : build

build: CFLAGS += $(BUILD_CFLAGS)
debug: CFLAGS += $(DEBUG_CFLAGS)


build: $(TARGET).so
debug: $(TARGET).so





$(TARGET).so: $(OBJECTS)
	$(call call_exec,${CC} $(LDFLAGS)  -o $@  $^, "LD      $@")
	




$(TARGET).a: $(OBJECTS)
	$(call call_exec,${AR} -rcs $@ $^,"AR      $@")


%.o:%.c ${HEADERS}
	$(call call_exec,${CC}  ${CFLAGS} ${CPU_CFLAGS} -c $< -o $@,"CC      $<")



clean:
	$(call call_exec,${RM} -rf $(OBJECTS) $(TARGET).so $(TARGET).a,"RM      $(OBJECTS)")
